# link-state simulation in go

For the final project of the computer network course at SBU University (spring 2021), we implemented a Link-state simulation in go.



You can read the full description in `todo.pdf`.



## Brief description

+ There is just one manager with all data (in YAML files)
+ It launches some router processes and connects to them with TCP
+ Router processes transfer their connectivity table so all of them get the whole network information
+ Manager command routers to transfer packets so the routing will be tested.



## run 

+ In order to run, you should run the manager executable 
+ But before that, the router should be compiled
+ There is a shell script responsible for compiling and running both



## configs

+ Manager read configs from Yaml files. They exist in the manager folder. fill free to change them
+ There should be just one `connection component` (all routers should connect to each other indirectly)
+ packets will be sent in order.



## logs

+ manager will print its own logs to stderr
+ Routers write logs to separate log files (with the name of their indices)
+ If a router fails before initializing the log file, it will send that failing log to the manager (the manager is watching the router processes stderr files)

+ `run.sh` will print all log files in order with one line between them, so you don't need to do anything in order to see logs.
