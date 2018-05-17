package main

import (
    "fmt"
    "time"
    "strconv"         //字符串转换
    "encoding/json"          //JSON包
    "github.com/gorilla/websocket"   //go的websocket实现
    "huobi/common"
    "net/http"
    "net/url"
    "huobi/conf"
)
type DataStruct struct {
    Open   float64 `json:"open"`
}
type ResStruct struct {
    Data     []DataStruct `json:"data"`
}

func main()  {
    dialer := websocket.Dialer{Proxy:func(*http.Request) (*url.URL, error) {
        return url.Parse("http://127.0.0.1:1080")
        },  //类型拨号器
}
    ws, _, err := dialer.Dial("wss://api.huobipro.com/ws", nil)  //拨号,返回连接对象,响应和错误
    if err != nil {
    	// handle error
    }

    for {        //因为是websocket,所以是不停的
        if _,p,err := ws.ReadMessage();err == nil {  //读取信息,这里是火币的规则
            res := common.UnGzip(p)            //解压数据
            //fmt.Println(string(res))           //输出字符串
            resMap := common.JsonDecodeByte(res)   //JSON解码
            if  v, ok := resMap["ping"];ok  {
                pingMap := make(map[string]interface{})
                pingMap["pong"] = v            //发送pong包,完成ping,pong
                pingParams := common.JsonEncodeMapToByte(pingMap)  //转成JSON
                if err := ws.WriteMessage(websocket.TextMessage, pingParams); err == nil {  //发送消息,TextMessage是整型,整数常量来标识两种数据消息类型
                    reqMap := new(common.ReqStruct)//创建结构体指针
                    reqMap.Id = strconv.Itoa(time.Now().Nanosecond())
                    reqMap.Req = conf.LtcTopic.KLineTopicDesc
                    reqBytes , err := json.Marshal(reqMap)
                    fmt.Println(string(reqBytes))
                    if err!=nil {
                        continue
                    }
                    if err := ws.WriteMessage(websocket.TextMessage,reqBytes); err == nil {  //发送ID REQ
                    }else{
                        fmt.Errorf("send req response error %s",err.Error())
                    }
                }else{
                    fmt.Errorf("huobi server ping client error %s",err.Error())
                    continue
                }
            }
            if  _, ok := resMap["rep"];ok  {//获取回应的API数据
                var resStruct ResStruct
                json.Unmarshal(res,&resStruct)
                //resStruct.Status
                //fmt.Printf("%+v",resStruct.Data)//输出内容
                for _,v:=range resStruct.Data{
                    fmt.Println(v.Open)
                }
            }
        }
    }
}
