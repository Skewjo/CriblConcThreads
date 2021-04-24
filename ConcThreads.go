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
const blockExecution=true;

func main(){
  if runtime.GOOS == "windows" {
    fmt.Print("OS: Windows\n")
  } 
  fmt.Print("Blocking calls:")
  if(blockExecution){
    fmt.Print("true\n")
  } else{
    fmt.Print("false\n")
  }
  fmt.Println()
  

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
  runTime := make([]time.Duration, len(binarySlice))

  for i := 0; i < len(binarySlice); i++{
    var concurrentGoRoutines = 0
    var maxConcurrentGoRoutines = 0
    
    var concThreads = 0
    var maxConcThreads = 0
    
    var timeStart = time.Now();

    var wg sync.WaitGroup
    for j:= 0; j < binarySlice[i]; j++{
      if(blockExecution){
        wg.Add(1)
      }
      //go worker(&wg, &concurrentGoRoutines)
      //In Go-Routine: Generating several random 64 bit floats with various delays
      go func(){
        concurrentGoRoutines++
        for h := 0; h < randIterations; h++{
          r := rand.New(rand.NewSource(42)) //probably need a "random" seed value...
          //r := rand.New(rand.NewSource(time.Now().UnixNano()) // Attempted, but failed with error: 
          _ = r
          time.Sleep(timeDelay * time.Second)
        }
        concurrentGoRoutines--
      }()
      if concurrentGoRoutines > maxConcurrentGoRoutines{
        maxConcurrentGoRoutines = concurrentGoRoutines
        //fmt.Printf("%d: %d\n", i, maxConcurrentGoRoutines)
      }
    }

    //I wanted to include the follow block (getting the thread count) inside of the GoRoutine, but this command call is too expensive and kills the program.
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
      if concThreads > maxConcThreads{
        maxConcThreads = concThreads
      }
    }
    var timeEnd = time.Now()
    var timeTotal = timeEnd.Sub(timeStart)

    concurrentRoutines[i] = maxConcurrentGoRoutines
    concurrentThreads[i] = maxConcThreads
    runTime[i] = timeTotal
  }

  //Print the results
  fmt.Println("\tTot. GR\t\tMax GR\t\tMax Threads\tRun Time")
  for i:= 0; i < len(binarySlice); i++{
    fmt.Printf("\t%d: ", binarySlice[i])  
    if(binarySlice[i] < 100000){
      fmt.Printf("\t")
    }
    fmt.Printf("\t%d", concurrentRoutines[i])
    fmt.Printf("\t\t%d", concurrentThreads[i])
    fmt.Print("\t\t")
    fmt.Print(runTime[i])
    fmt.Print("\n")
  }
}
