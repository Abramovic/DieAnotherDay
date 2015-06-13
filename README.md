# DieAnotherDay

A simple tool that does it's best to avoid processes from getting killed. 

## Using the Tool 

#### Basic Example

To build:

````golang
go build .
````

To run:

````
./dieanotherday "my-command" 
````

#### Daemonize 

To have this tool run in the background add the flag "-d=true" before your process

````
./dieanotherday -d=true "my-command"
````

#### Process Args 

If you want to give your process some additional arguments just add them at the end 

````
./dieanotherday "my-command" "-arg1=test"
````

#### Customize Timeouts

By default the tool will check every second to see if your process is still up. If you want to change the timing you can include the "t" argument. 

````
./dieanotherday -t=30 "my-command" "-arg1=test"
````

## Managing the Process

If you want to view the output of your child process you can view the file ````dieanotherday.log````
