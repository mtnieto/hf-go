# Proof of concept - Hyperledger Fabric Go SDK

## Requisitos ##
Usar el SDK de fabric para Go en una versión concreta, los cambios entre commits son muy diferentes y puede que el cliente deje de funcionar.

```
go get -u github.com/hyperledger/fabric-sdk-go && \
cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && \
git checkout 614551a752802488988921a730b172dada7def1d
```

## Configuración ##

Hay que añadir los paths de los certificados y las urls donde se encuentran el nodo sobre el cual se quiera operar.

Por otro lado a la hora de debug, y sacar trazas hay que modificar lo siguiente:

```
  logging:
    level: debug
```

La configuración del yaml tira contra la configuración de ejemplo de github:

https://github.com/mtnieto/hf-2orgs

NOTA: Puede usar cualquiera. Siempre y cuando se estblezcan bien los paths.

## main.go ##

Como prerrequisitos hay que tener la red levantada, los canales a usar creados y los peer que hayan hecho join en el mismo, Sino se producirá un error de unknown xxxxMSP.

Lo que realiza esta poc es la instalación de un contrato y su posterior instanciación (despligue en la red de mismo).


