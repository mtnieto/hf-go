# Proof of concept - Hyperledger Fabric Go SDK


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


