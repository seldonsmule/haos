package main

import (

  "fmt"
  "flag"
  "os"
  "strings"
  "time"
  //"/Users/rxe789/egc_NetworkDrive/dev/gocode/src/github.com/seldonsmule/unifi"
  //"/Users/rxe789/dev/golang/src/github.com/seldonsmule/haos"
  "github.com/seldonsmule/haos"
  "github.com/seldonsmule/simpleconffile"
  "github.com/seldonsmule/logmsg"

)

type Configuration struct {

  API_KEY string // Key from HA
  ConfFilename string // name of conf file
  Host string // name of the HA server
  Encrypted bool

}

const COMPILE_IN_KEY = "example key 9999"

var gMyConf Configuration

func readconf(confFile string, printstd bool) bool{

  simple := simpleconffile.New(COMPILE_IN_KEY, confFile)
  
  if(!simple.ReadConf(&gMyConf)){
    msg := fmt.Sprintln("Error reading conf file: ", confFile)
    logmsg.Print(logmsg.Warning, msg)
    return false
  }
  
  if(gMyConf.Encrypted){    
    gMyConf.API_KEY = simple.DecryptString(gMyConf.API_KEY)
  }


  if(printstd){

    fmt.Printf("Encrypted [%v]\n", gMyConf.Encrypted)
    fmt.Printf("API_KEY [%v]\n", gMyConf.API_KEY)
    fmt.Printf("Host [%v]\n", gMyConf.Host)
    fmt.Printf("ConfFilename [%v]\n", gMyConf.ConfFilename)

  }

  return true

}

func help(){

  fmt.Println("cmd not found")

  flag.PrintDefaults()

  fmt.Println("cmds:")
  fmt.Println("      setconf - Setup Conf file")
  fmt.Println("        -apikey HomeAssistant API Key")
  fmt.Println("        -host URL of HomeAssistant host")
  fmt.Println("        -conffile name of conffile (.ha.conf default)")
  fmt.Println()
  fmt.Println("      readconf - Display Conf file info")
  fmt.Println()
  fmt.Println("      get_structs - get some response structs and save off")
  fmt.Println()
  fmt.Println("      get_entity_state - Display Conf file info")
  fmt.Println("        -entity of HomeAssistant entity")
  fmt.Println()
  fmt.Println("      set_entity_state - Sets a known entity's state")
  fmt.Println("        -entity of HomeAssistant entity")
  fmt.Println("        -state new state")
  fmt.Println()
  fmt.Println("      set_entity_state_on - Sets a known entity's state to on")
  fmt.Println()
  fmt.Println("      set_entity_state_off - Sets a known entity's state to off")




}

func main(){

  cmdPtr := flag.String("cmd", "help", "Command to run")
  apikeyPtr := flag.String("apikey", "notset", "HomeAssistant API KEY")
  hostPtr := flag.String("host", "notset", "URL of Host system")
  entityPtr := flag.String("entity", "notset", "Name of an HA entity")
  statePtr := flag.String("state", "notset", "State to set of an HA entity")
  confPtr := flag.String("conffile", ".ha.conf", "config file name")
  bdebugPtr := flag.Bool("debug", false, "If true, do debug magic")

  flag.Parse()

  fmt.Printf("cmd=%s\n", *cmdPtr)

  logmsg.SetLogFile("example.log");

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "apikeyPtr = ", *apikeyPtr)
  logmsg.Print(logmsg.Info, "hostPtr = ", *hostPtr)
  logmsg.Print(logmsg.Info, "entityPtr = ", *entityPtr)
  logmsg.Print(logmsg.Info, "statePtr = ", *statePtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "bdebugPtr = ", *bdebugPtr)


  fmt.Println("starting");

  readconf(*confPtr, false);



  ha := haos.New(gMyConf.Host)
  //ha.SetApiKey(gMyConf.API_KEY)
  ha.SetToken(gMyConf.API_KEY)

  ha.Dump()

  switch *cmdPtr {

    case "readconf":
      fmt.Println("Reading Conf File")
      readconf(*confPtr, true)

    case "setconf":

      readconf(*confPtr, false); // ignore errors

      fmt.Println("Setting conf file")

      simple := simpleconffile.New(COMPILE_IN_KEY, *confPtr);

      gMyConf.Encrypted = true

      if(strings.Compare(*apikeyPtr, "notset") != 0){
       gMyConf.API_KEY = simple.EncryptString(*apikeyPtr)
      }else{
       gMyConf.API_KEY = simple.EncryptString(gMyConf.API_KEY)
      } 

      if(strings.Compare(*hostPtr, "notset") != 0){
        gMyConf.Host = *hostPtr
      }

      gMyConf.ConfFilename = *confPtr


      simple.SaveConf(gMyConf)
 
/*
    case "listsiteids":
     fmt.Printf("Site IDs for Host[%s]\n", gMyConf.Host)
     //un.SetApiKey(gMyConf.API_KEY)
     un.ListSitesIDs()
*/

    case "get_entity_state":
      if(ha.GetEntityState(*entityPtr, false)){
        fmt.Printf("Get EntityState[%s] worked\n", *entityPtr)
        fmt.Printf("Current State[%s]\n", ha.EntityState.State)
        loc, _ := time.LoadLocation("America/New_York")
        lastchg := ha.EntityState.LastChanged.In(loc)
        fmt.Println("LastChanged: " + lastchg.String())
        lastrpt := ha.EntityState.LastReported.In(loc)
        fmt.Println("LastReported: " + lastrpt.String())
        lastup := ha.EntityState.LastUpdated.In(loc)
        fmt.Println("LastUpdated: " + lastup.String())
      }else{
        fmt.Printf("Get EntityState[%s] failed\n", *entityPtr)
      }

    case "set_entity_state_on":
      //1st make sure the entity already exist
      if(ha.GetEntityState(*entityPtr, false)){
 
        ha.SetEntityStateOn(*entityPtr)

      }else{
        fmt.Printf("[%s] does not exist\n", *entityPtr)
      }

    case "set_entity_state_off":
      //1st make sure the entity already exist
      if(ha.GetEntityState(*entityPtr, false)){
 
        ha.SetEntityStateOff(*entityPtr)

      }else{
        fmt.Printf("[%s] does not exist\n", *entityPtr)
      }

    case "set_entity_state":

      if( *statePtr == "notset"){
        fmt.Println("-state not passed in - exiting")
        break
      }
      //1st make sure the entity already exist
      if(ha.GetEntityState(*entityPtr, false)){
 
        ha.SetEntityState(*entityPtr, *statePtr)

      }else{
        fmt.Printf("[%s] does not exist\n", *entityPtr)
      }
        


    case "get_structs":
      ha.GetStates(true)
      ha.GetEntityState("update.home_assistant_supervisor_update",true)
      fmt.Printf("Number of states [%d]\n", len(ha.States))


      for i:=0; i < len(ha.States); i++ {

        fmt.Printf("[%d] - [%s] State[%s]\n", i, ha.States[i].EntityID,
                                                ha.States[i].State)

      }


    

    default:
      help()
      os.Exit(2)

  }


  
}
