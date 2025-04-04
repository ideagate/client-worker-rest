# Client Worker REST

# Layer Structure
```mermaid
graph BT
    coreCache[Core/Cache]
    adapterController[Adapter Controller]
    usecaseEndpoint[Use Case Endpoint]
    usecaseHandler[Use Case Handler]
    router[Router]
    
    usecaseHandler --> router
    usecaseEndpoint --> usecaseHandler
    coreCache --> usecaseEndpoint
    adapterController --> usecaseEndpoint
```