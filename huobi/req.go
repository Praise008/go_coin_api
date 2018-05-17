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
func main()  {
    dialer := websocket.Dialer{Proxy:func(*http.Request) (*url.URL, error) {
        return url.Parse("http://127.0.0.1:1080")
        },  //类型拨号器
}
    ws, _, err := dialer.Dial("wss://api.huobipro.com/ws", nil)  //拨号,返回连接对象,响应和错误
    if err != nil {
    	// handle error
    }
    reqMap := new(common.SubStruct)//创建结构体指针
    reqMap.Id = strconv.Itoa(time.Now().Nanosecond())
    reqMap.Sub = conf.LtcTopic.KLineTopicDesc
    reqBytes , err := json.Marshal(reqMap)
    //fmt.Println(string(reqBytes))
    errr := ws.WriteMessage(websocket.TextMessage,reqBytes)
    if errr != nil{
        fmt.Errorf(err.Error())
    }
    
    for {        //因为是websocket,所以是不停的
        if _,p,err := ws.ReadMessage();err == nil {
            res := common.UnGzip(p)            //解压数据
            //fmt.Println(string(res))           //输出字符串
            resMap := common.JsonDecodeByte(res)   //JSON解码
            if  v, ok := resMap["ping"];ok  {
                pingMap := make(map[string]interface{})
                pingMap["pong"] = v            //发送pong包,完成ping,pong
                pingParams := common.JsonEncodeMapToByte(pingMap)  //转成JSON
                err := ws.WriteMessage(websocket.TextMessage, pingParams)
                if err != nil{
                    fmt.Errorf(err.Error())
                }
            }
            if  vv, okk := resMap["tick"];okk  {
                bch_values := vv.(map[string]interface{})         
                fmt.Println(bch_values["open"])
            }
            
        }
    }
}
