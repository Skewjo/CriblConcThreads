# CriblConcThreads

## Thread Use Comparison Between Windows and WSL
While the run times are similar (in this excercise at least), it appears that Golang takes a much different approach in regards to thread usage between the different OS:
| Windows | Linux |
| ------------- | ------------- |
| ![real_results_windows](https://user-images.githubusercontent.com/4118039/116083614-6c5d6d00-a662-11eb-86ab-496e78c5a720.PNG) | ![real_results_linux](https://user-images.githubusercontent.com/4118039/116083617-6ebfc700-a662-11eb-895d-b0bedcf59187.PNG) |

## Little Optimization for Fixed Seeds
Somewhat surprisingly, Go performs less optimization than I had expected for using the same fixed random seed in all 100000+ Go-routines:
| Windows | Linux |
| ------------- | ------------- |
| ![no_optimization_for_fixed_seeds](https://user-images.githubusercontent.com/4118039/116079916-0969d700-a65e-11eb-9369-813708a4c3d2.PNG) | ![optimization for fixed seeds](https://user-images.githubusercontent.com/4118039/116081756-40d98300-a660-11eb-92bc-90cf58a48116.png) |

## Wait Groups are expensive
The act of "```add()```'ing" to a wait group appears to be somewhat expensive. This was a preliminary version of the program where I had accidentally omitted the ```defer()``` and ```wait()``` functions necessary to make blocking calls function as intended. As you can see, the execution time shoots up by ~25% and the maximum number of concurrent Go-routines skyrockets by 400%:
![defer() and wait() omitted](https://user-images.githubusercontent.com/4118039/116082178-ba717100-a660-11eb-9cbf-9439e0e2a8ca.png)

## Cranked to 11
I decided to up the number of Go-routines created up from 100,000+ to 1,000,000+, and at some point, both Windows and Linux fall off the rails:

| Windows | Linux |
| ------------- | ------------- |
| ![cranked_to_11](https://user-images.githubusercontent.com/4118039/116080164-551c8080-a65e-11eb-9ba4-4ab82ff018a8.PNG)  | ![linux_1_mil](https://user-images.githubusercontent.com/4118039/116080238-70878b80-a65e-11eb-854d-bf2f15e6e4a0.PNG) |




## System Info
Here's the system info for the machine I've run these tests on:
![sys_info](https://user-images.githubusercontent.com/4118039/116079820-e4756400-a65d-11eb-8ff9-7a77797d0f29.PNG)
