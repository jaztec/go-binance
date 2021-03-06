package model

import "encoding/json"

// StreamData holds the upmost fields of a stream package. The data can represent any
// described package.
type StreamData struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"executionReport","E":1619561096919,"s":"DOGEEUR","c":"web_aad35285a92e4de3a614dd3a2f0dcaad","S":"SELL","o":"TAKE_PROFIT_LIMIT","f":"GTC","q":"140.30000000","p":"0.23210000","P":"0.23200000","F":"0.00000000","g":-1,"C":"web_dcd5475d96544df8980bbf0189ffbeb4","x":"CANCELED","X":"CANCELED","r":"NONE","i":51119693,"l":"0.00000000","z":"0.00000000","L":"0.00000000","n":"0","N":null,"T":1619561096919,"t":-1,"I":109777202,"w":false,"m":false,"M":false,"O":1619551088231,"Z":"0.00000000","Y":"0.00000000","Q":"0.00000000"}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"outboundAccountPosition","E":1619561096919,"u":1619561096919,"B":[{"a":"BNB","f":"0.02608043","l":"0.00000000"},{"a":"DOGE","f":"140.32390000","l":"0.00000000"},{"a":"EUR","f":"6.60152400","l":"0.00000000"}]}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"executionReport","E":1619561115925,"s":"DOGEEUR","c":"web_c935a242aecc48c6b4a262a48a734951","S":"SELL","o":"LIMIT","f":"GTC","q":"140.30000000","p":"0.22700000","P":"0.00000000","F":"0.00000000","g":-1,"C":"","x":"NEW","X":"NEW","r":"NONE","i":51202991,"l":"0.00000000","z":"0.00000000","L":"0.00000000","n":"0","N":null,"T":1619561115924,"t":-1,"I":109777512,"w":true,"m":false,"M":false,"O":1619561115924,"Z":"0.00000000","Y":"0.00000000","Q":"0.00000000"}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"outboundAccountPosition","E":1619561115925,"u":1619561115924,"B":[{"a":"BNB","f":"0.02608043","l":"0.00000000"},{"a":"DOGE","f":"0.02390000","l":"140.30000000"},{"a":"EUR","f":"6.60152400","l":"0.00000000"}]}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"executionReport","E":1619561123516,"s":"DOGEEUR","c":"web_c935a242aecc48c6b4a262a48a734951","S":"SELL","o":"LIMIT","f":"GTC","q":"140.30000000","p":"0.22700000","P":"0.00000000","F":"0.00000000","g":-1,"C":"","x":"TRADE","X":"PARTIALLY_FILLED","r":"NONE","i":51202991,"l":"11.40000000","z":"11.40000000","L":"0.22700000","n":"0.00000411","N":"BNB","T":1619561123516,"t":8351848,"I":109777661,"w":false,"m":true,"M":true,"O":1619561115924,"Z":"2.58780000","Y":"2.58780000","Q":"0.00000000"}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"outboundAccountPosition","E":1619561123516,"u":1619561123516,"B":[{"a":"BNB","f":"0.02607632","l":"0.00000000"},{"a":"DOGE","f":"0.02390000","l":"128.90000000"},{"a":"EUR","f":"9.18932400","l":"0.00000000"}]}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"executionReport","E":1619561125015,"s":"DOGEEUR","c":"web_c935a242aecc48c6b4a262a48a734951","S":"SELL","o":"LIMIT","f":"GTC","q":"140.30000000","p":"0.22700000","P":"0.00000000","F":"0.00000000","g":-1,"C":"","x":"TRADE","X":"FILLED","r":"NONE","i":51202991,"l":"128.90000000","z":"140.30000000","L":"0.22700000","n":"0.00004650","N":"BNB","T":1619561125014,"t":8351849,"I":109777693,"w":false,"m":true,"M":true,"O":1619561115924,"Z":"31.84810000","Y":"29.26030000","Q":"0.00000000"}}
//{"stream":"kQcTHLQWj24fRzmxeIf9zqepIH7f5Sai3rblhbGSQpQlZFPtVdmQK3j4J91m","data":{"e":"outboundAccountPosition","E":1619561125015,"u":1619561125014,"B":[{"a":"BNB","f":"0.02602982","l":"0.00000000"},{"a":"DOGE","f":"0.02390000","l":"0.00000000"},{"a":"EUR","f":"38.44962400","l":"0.00000000"}]}}
