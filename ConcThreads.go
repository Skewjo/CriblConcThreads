package main

import(
  "fmt"
  "runtime"
  "math/rand"
  "time"
  "sync"
  "os"
  "os/exec"
  "strconv"
  "log"
  "syscall"
  "strings"
)

//These 4 constants are used to the determine the complexity of the GoRoutine. 
const timeDelay=1;
const randIterations=3;
const blockExecution=false;
const fixedSeed=true;

//TODO: Make functions for:
// goRoutineWorker

//TODO: Make it work on Linux
//TODO: Get better grasp on GoRoutine sync calls, defer, & sleep
//TODO: Show more system info

func main(){
  printSysAndRTInfo()
  
  //The following block creates a slice of numbers from 2^0 up to 2^17
  binarySlice := make([]int, 1)
  binarySlice[0] = 1
  for i:= 0; binarySlice[i] < 1000000; i++{
    binarySlice = append(binarySlice, binarySlice[i]*2)
    if binarySlice[i] > 1000000{
      break
    }
  }

  concurrentRoutines := make([]int, len(binarySlice))
  concurrentThreads := make([]int, len(binarySlice))
  additionalThreads := make([]int, len(binarySlice))
  runTime := make([]time.Duration, len(binarySlice))
  
  origThreadCount := 0
  if runtime.GOOS == "windows"{
    origThreadCount = getThreadCountWindows()
  } else if runtime.GOOS == "linux"{
    origThreadCount = getThreadCountLinux()
  }
  
  fmt.Printf("Thread count before work: %d\n", origThreadCount)

  timeStartTotal := time.Now()

  //For each value from 2^0 to 2^n
  for i := 0; i < len(binarySlice); i++{
    concurrentGoRoutines := 0
    maxConcurrentGoRoutines := 0
    
    concThreads := 0
    maxConcThreads := 0
    
    timeStart := time.Now()

    var wg sync.WaitGroup

    //Create specified number of GoRoutines
    for j:= 0; j < binarySlice[i]; j++{
      if(blockExecution){
        wg.Add(1)
      }
      //go worker(&wg, &concurrentGoRoutines)
      //In Go-Routine: Generating several random 64 bit floats with various delays
      go func(){
        concurrentGoRoutines++
        if(blockExecution){
          defer wg.Done() //"On return, notify the WaitGroup that we're done" 
        }
        
        for h := 0; h < randIterations; h++{
          var r *rand.Rand
          if(fixedSeed){
            r = rand.New(rand.NewSource(42))
          } else{
            r = rand.New(rand.NewSource(time.Now().UnixNano())) 
          }
          _ = r //Throw the value away
          time.Sleep(timeDelay * time.Second)
        }
        concurrentGoRoutines--
      }()
      
      if concurrentGoRoutines > maxConcurrentGoRoutines{
        maxConcurrentGoRoutines = concurrentGoRoutines
      }
    }
    if(blockExecution){
      wg.Wait()
    }

    //I wanted to include the following block inside of the GoRoutine, but this command call(on Windows at least) is too expensive and locks up the program. Ultimately I don't believe this matters because I found that the threads take quite a bit of time to be cleaned up (on the order of seconds).
    if runtime.GOOS == "windows" {
      concThreads = getThreadCountWindows()
      if concThreads > maxConcThreads{
        maxConcThreads = concThreads
      }
    }else if runtime.GOOS == "linux" {
      concThreads = getThreadCountLinux()
      if concThreads > maxConcThreads{
        maxConcThreads = concThreads
      }
    }
    var timeEnd = time.Now()
    var timeTotal = timeEnd.Sub(timeStart)

    concurrentRoutines[i] = maxConcurrentGoRoutines
    concurrentThreads[i] = maxConcThreads
    additionalThreads[i] = maxConcThreads - origThreadCount
    runTime[i] = timeTotal
  }
  
  var timeEndTotal = time.Now()
  var timeTotalTotal = timeEndTotal.Sub(timeStartTotal)

  fmt.Print("Total Run Time: ", timeTotalTotal, "\n")
  printResults(binarySlice, concurrentRoutines, concurrentThreads, additionalThreads, runTime)
}

func printSysAndRTInfo(){
  fmt.Println("System info:")
  fmt.Printf("\tGo version: %s\n", runtime.Version())
  fmt.Printf("\tOS: %s\n", runtime.GOOS)
  fmt.Println()

  fmt.Println("Options:")
  fmt.Printf("\tTime Delay: %d seconds\n", timeDelay)
  fmt.Printf("\tRandIterations: %d\n", randIterations)
  fmt.Print("\tBlocking calls: ")
  if(blockExecution){
    fmt.Print("true\n")
  } else{
    fmt.Print("false\n")
  }
  fmt.Print("\tFixed Seed: ")
  if(fixedSeed){
    fmt.Print("true\n")
  } else{
    fmt.Print("false\n")
  }
  fmt.Println()
}

func getThreadCountLinux() int{
  var concThreads = 0
  if runtime.GOOS == "linux"{
    stdOut, err := exec.Command("ps", "huH", "p", strconv.Itoa(os.Getpid())).Output()
    if err != nil {
      fmt.Sprintf("Failed to execute command: %s", err)
    }
    concThreads = len(strings.Split(string(stdOut), "\n")) - 1
  }
  return concThreads;
}

func getThreadCountWindows() int{
  var concThreads = 0
  if runtime.GOOS == "windows" {
    cmd := exec.Command("./GetThreadCount.exe", strconv.Itoa(os.Getpid()))
    if err := cmd.Start(); err != nil {
      log.Fatalf("cmd.Start: %v", err)
    }
    if err := cmd.Wait(); err != nil {
      if exiterr, ok := err.(*exec.ExitError); ok {
        if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
          concThreads = status.ExitStatus()
        }
      } else {
        log.Fatalf("cmd.Wait: %v", err)
      }
    }
  }
  return concThreads;
}

func printResults(binarySlice []int, concurrentRoutines []int, concurrentThreads []int, additionalThreads []int, runTime []time.Duration){
  fmt.Println("\tTot. GR\t\tMax Conc. GR\tMax Threads\tAddl. Threads\tRun Time")
  for i:= 0; i < len(binarySlice); i++{
    fmt.Printf("\t%d: ", binarySlice[i])  
    if(binarySlice[i] < 100000){
      fmt.Printf("\t")
    }
    fmt.Printf("\t%d", concurrentRoutines[i])
    fmt.Printf("\t\t%d", concurrentThreads[i])
    fmt.Printf("\t\t%d", additionalThreads[i])
    fmt.Print("\t\t", runTime[i], "\n")
  }
}