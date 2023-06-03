# smart-plug
Smart Plug is a microservice that turns on/off and monitors smart plugs "Shelly" depending on conditions that are set by the user.

The server uses GraphQL for its API.

## To add a new plug for monitoring use the following query
```
mutation CreatePlug {
    createPlug(input: {
        ipAddress: "<PLUG_IP>"
        name: "<PLUG_NAME>"
        powerToTurnOff: <MINIMAL_POWER_TO_KEEP_CHARGING>
    }) {
        id
        ipAddress
        name
        powerToTurnOff
        createdAt
    }
}
```


## To list all plugs use
```
query ListPlugs {
    listPlugs(page: 1, perPage: 50) {
        id
        ipAddress
        name
        powerToTurnOff
        createdAt
    }
}
```
