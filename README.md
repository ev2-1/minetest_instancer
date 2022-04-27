# minetest_instancer

Docker based minetest auto server controlled over telnet written in go

A client for this is [my telnetClient](https://github.com/ev2-1/telnetClient) library

Volumes:
- `/var/run/docker.sock` - mount docker socket here
- `/mount/worlds` - a place where your worlds folder has to be mounted to create worlds/folders

Configuration: (config.json)
- `SrvNetwork` - the network the containers will be added to
- `MinetestContainer` - the container minetest servers will be created from (aka ./minetest git submodules)
- `WorldPath` - the **ABSOLUTE** path on the **HOST MACHINE** where the worlds should be saved (**HAS** to be the same as whats mounted on `/mount/worlds`) 
- `ConfigPath` - the **ABSOLUTE** path on the **HOST MACHINE** where minetest configuration files are mounted from
- `GamePath` - the **ABSOLUTE** path on the **HOST MACHINE** where minetest Games will be loaded from

Telnet c&c interface (port 8888):

`get_ip`

> Usage: `get_ip <container>`
>
> Return value: `OK, <ip>`
>
> Gets the ip address of docker container in network


`srv_create`

> Usage: `srv_create <name> <world> <game> [net]`
>
> Return value: `OK, <containerID>`
>
> creates server based on container**name**, world-name, gameid and optionally the network the container will be connected to (otherwise won't be connected automatically)
>
> use `default` for usage of configured `SrvNetwork`

`srv_connect`

> Usage: `srv_connect <container> <network>`
>
> Return value: `OK`
>
> connects existing container to network
>
> use `default` for usage of configured `SrvNetwork`

`srv_delete`

> Usage: `srv_delete <container>`
>
> Return value: `OK`
>
> delets server based on id

`srv_state`

> Usage: `srv_state <container>`
>
> Return value: 
>
> - `OK, false` - container does not exist
>
> - `OK, <status>, <Running>`
>
> gets state of container,
>
> state - either `created`, `running`, `restarting`, `exited`, `paused`, `dead` [detailed explenation](https://www.baeldung.com/ops/docker-container-states)
>
> Running - Either `Running` or `Stopped`

`servers`

> Usage: `servers`
>
> Return value:
>
> list of servers seperated by `;` in format:
>
> `<id> - <name>|<world>|<game>|<config>|<full_id>|<net>|<ip>|<ready>`
>
> ready - true if server got `instancer/ready`

`srv_start`

> Usage: `srv_start <container>`
>
> Return value: `OK`
>
> start server based on id

Errors:

> `err, <error>` - where error is a error message
>
> `err` - undefined error
