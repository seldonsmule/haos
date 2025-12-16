package haos

import (

  "fmt"
  "encoding/json"
//  "flag"
//  "os"
//  "strings"
//  //"/Users/rxe789/egc_NetworkDrive/dev/gocode/src/github.com/seldonsmule/unifi"
//  "github.com/seldonsmule/unifi"
  "github.com/seldonsmule/restapi"
  "github.com/seldonsmule/logmsg"

)

type Haos struct {

  sToken string // From the HOAS console
  sBaseEndpoint string // Base endpoint

  States HaosStates // List of entity states 
  EntityState HaosEntityState // A single entit - cheated and grab a known one


}

func New(URL string) *Haos {

  ha := new(Haos)

  ha.sBaseEndpoint = URL + "/api"

  return(ha)
}

func (pHa *Haos) Dump(){

  fmt.Printf("Dump of Haos\n")
  fmt.Printf("Haos.sBaseEndpoint[%s]\n", pHa.sBaseEndpoint)
  //fmt.Printf("Haos.sToken[%s]\n", pHa.sToken)

}

func (pHa *Haos) SetToken(sToken string) bool{

  pHa.sToken = sToken

  return true
}

func (pHa *Haos) SetEntityStateOn(sEntity string) bool{
  return(pHa.SetEntityState(sEntity, "on"))
}

func (pHa *Haos) SetEntityStateOff(sEntity string) bool{
  return(pHa.SetEntityState(sEntity, "off"))
}

func (pHa *Haos) SetEntityState(sEntity string, sNewState string) bool{

  msg := fmt.Sprintf("POST state change for [%s] to  [%s]", sEntity, sNewState)
  logmsg.Print(logmsg.Info, msg)

  endpointname := pHa.sBaseEndpoint + "/states/" + sEntity
  
  logmsg.Print(logmsg.Info, "SetEntityState: " + endpointname)
  jsondata := fmt.Sprintf("{\"state\" : \"%s\" }", sNewState)

  logmsg.Print(logmsg.Info, "SetEntityState: " + jsondata)

  r := restapi.NewPost("setstates", endpointname)

  r.SetBearerAccessToken(pHa.sToken)

  r.SetPostJson(jsondata)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

//  r.Dump()

//  r.DebugOn()


  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  logmsg.Print(logmsg.Info,r.BodyString)


  return true
}

func (pHa *Haos) AutomationTrigger(sAutomation string) bool{

  msg := fmt.Sprintf("POST automation trigger for [%s]", sAutomation)
  logmsg.Print(logmsg.Info, msg)

  endpointname := pHa.sBaseEndpoint + "/services/automation/trigger"
  
  logmsg.Print(logmsg.Info, "URL: " + endpointname)
  
  jsondata := fmt.Sprintf("{\"entity_id\" : \"automation.%s\" }", sAutomation)

  logmsg.Print(logmsg.Info, "jsondata: " + jsondata)

  r := restapi.NewPost("automationtrigger", endpointname)

  r.SetBearerAccessToken(pHa.sToken)

  r.SetPostJson(jsondata)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

  //r.Dump()

  //r.DebugOn()


  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  logmsg.Print(logmsg.Info,r.BodyString)


  return true
}

func (pHa *Haos) GetEntityState(sEntity string, bSave bool) bool{

  endpointname := pHa.sBaseEndpoint + "/states/" + sEntity

  logmsg.Print(logmsg.Info, "GetEntityState: " + endpointname)

  r := restapi.NewGet("getstates", endpointname)

  r.SetBearerAccessToken(pHa.sToken)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

//  r.Dump()

//  r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  // cheating and saving off
  if(bSave){
    r.SaveResponseBody("haos_entitystate", "HaosEntityState", false)
  }

  json.Unmarshal(r.BodyBytes, &pHa.EntityState)

  return true

}

func (pHa *Haos) GetStates(bSave bool) bool{

  endpointname := pHa.sBaseEndpoint + "/states"

  r := restapi.NewGet("getstates", endpointname)

  r.SetBearerAccessToken(pHa.sToken)

  restapi.TurnOffCertValidation()

  r.JsonOnly()

//  r.Dump()

//  r.DebugOn()

  if(!r.Send()){
    msg := fmt.Sprintf("Error getting [%s]\n", endpointname)
    //fmt.Printf("Error sending: %s\n", msg)
    logmsg.Print(logmsg.Error, msg)
    return false
  }

  if(bSave){
    r.SaveResponseBody("haos_states", "HaosStates", false)
  }

  // cheating and saving off

  json.Unmarshal(r.BodyBytes, &pHa.States)

  return true

}

