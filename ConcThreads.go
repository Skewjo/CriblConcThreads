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
)

const timeDelay=10;
const randIterations=3;
const blockExecution=false;
const fixedSeed=false;

//TODO: Make functions for:
// 2. populateSlice
// 3. goRoutineWorker

//TODO: Make it work on Linux
//TODO: Get better grasp on GoRoutine sync calls, defer, & sleep
//TODO: Show more system info

func main(){
  printSysAndRTInfo()
  
  //The following block creates a slice of numbers from 2^0 up to 2^17
  binarySlice := make([]int, 1)
  binarySlice[0] = 1
  for i:= 0; binarySlice[i] < 100000; i++{
    binarySlice = append(binarySlice, binarySlice[i]*2)
    if binarySlice[i] > 100000{
      break
    }
  }

  concurrentRoutines := make([]int, len(binarySlice))
  concurrentThreads := make([]int, len(binarySlice))
  additionalThreads := make([]int, len(binarySlice))
  runTime := make([]time.Duration, len(binarySlice))
  
  origThreadCount := getThreadCountWindows()
  fmt.Printf("Thread count before work: %d\n", origThreadCount)

  //For each value from 2^0 to 2^n
  for i := 0; i < len(binarySlice); i++{
    var concurrentGoRoutines = 0
    var maxConcurrentGoRoutines = 0
    
    var concThreads = 0
    var maxConcThreads = 0
    
    var timeStart = time.Now();

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

    //I wanted to include the follow block inside of the GoRoutine, but this command call(on Windows at least) is too expensive and locks up the program. Ultimately I don't believe this matters because I found that the threads take quite a bit of time to be cleaned up (on the order of seconds).
    if runtime.GOOS == "windows" {
      concThreads = getThreadCountWindows()
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

  //Print the results
  printResults(binarySlice, concurrentRoutines, concurrentThreads, additionalThreads, runTime)
}

func printSysAndRTInfo(){
  fmt.Printf("Go version: %s\n", runtime.Version())
  fmt.Printf("OS: %s\n", runtime.GOOS)
  fmt.Print("Blocking calls: ")
  if(blockExecution){
    fmt.Print("true\n")
  } else{
    fmt.Print("false\n")
  }
  fmt.Print("Fixed Seed: ")
  if(fixedSeed){
    fmt.Print("true\n")
  } else{
    fmt.Print("false\n")
  }
  fmt.Println()
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
  fmt.Println("\tTot. GR\t\tMax GR\t\tMax Threads\tAddl. Threads\tRun Time")
  for i:= 0; i < len(binarySlice); i++{
    fmt.Printf("\t%d: ", binarySlice[i])  
    if(binarySlice[i] < 100000){
      fmt.Printf("\t")
    }
    fmt.Printf("\t%d", concurrentRoutines[i])
    fmt.Printf("\t\t%d", concurrentThreads[i])
    fmt.Printf("\t\t%d", additionalThreads[i])
    fmt.Print("\t\t")
    fmt.Print(runTime[i])
    fmt.Print("\n")
  }
}