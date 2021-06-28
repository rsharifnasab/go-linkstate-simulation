# link-state simulation in go

for final project of computer-network course in SBU university (spring 2021) we implemented a Link-state simulation in go.



you can read full description in `todo.pdf`



## brief description

+ there is just one manager with all data's (in yaml files)
+ it launch some router processes and the connect to them with tcp
+ router processes transfer their connectivity table so all of them get whole network information
+ manager command routers to transfer packets, so the routing will be tested.



## run 

+ in order to run, you should run manager executable 
+ but before that, router should be compiled
+ there is a shell script responsible for compiling and running both



## configs

+ manager read configs from Yaml files. there are exists in manager folder. fill free to change them
+ there should be just one `connection component` (all routers should connect to each other indirectly)
+ packets will send in order.



## logs

+ manager will print its own logs to stderr
+ routers write logs to separate log files (with name of their indices)
+ if a router fail before initializing the log file, it will send that failing log to manager (manager is watching router processes stderr files)

+ `run.sh` will print all log files in order with one line between them, so you don't need to do anything in order to see logs.





