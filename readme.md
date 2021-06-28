# link-state simulation in go

for final project of computer-network course in SBU university (spring 2021) we implemented a Link-state simulation in go.



you can read full description in `todo.pdf`



## brief description

+ there is just one manager with all data's (in yaml files)
+ it launch some router processes and the connect to them with tcp
+ router processes transfer their connectivity table so all of them get whole network information
+ manager command routers to transfer packets, so the routing will be tested.

