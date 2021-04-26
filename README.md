# CriblConcThreads

Somewhat surprisingly, Go does not perform any optimization for using the same fixed random seed in all 100000+ Go-routines:
![no_optimization_for_fixed_seeds](https://user-images.githubusercontent.com/4118039/116079916-0969d700-a65e-11eb-9369-813708a4c3d2.PNG)

I decided to up the number of Go-routines created up from 100,000+ to 1,000,000+, and at some point, both Windows and Linux fall off the rails:
![cranked_to_11](https://user-images.githubusercontent.com/4118039/116080164-551c8080-a65e-11eb-9ba4-4ab82ff018a8.PNG)
![linux_1_mil](https://user-images.githubusercontent.com/4118039/116080238-70878b80-a65e-11eb-854d-bf2f15e6e4a0.PNG)

## System Info
Here's the system info for the machine I've run these tests on:
![sys_info](https://user-images.githubusercontent.com/4118039/116079820-e4756400-a65d-11eb-8ff9-7a77797d0f29.PNG)
