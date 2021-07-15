#!/bin/bash

# Script segue os seguintes passos:
# - Atualizar o repsitorio local
# - Fazer o build do container 
# - Startar o serviço do "GoDebut"

clear

HASH="##############################################"

echo $HASH
echo "Iniciando PULL do repositorio..."
echo $HASH
echo " "

git pull

echo $HASH
echo "Iniciando BUILD da imagem..."
echo $HASH
echo " "

docker build -t ffelixneto/godebut:latest .

echo $HASH
echo "Iniciando o serviço GoDebut..."
echo $HASH
echo " "

docker stack deploy godebut -c docker-compose.yml

sleep 3

echo $HASH
echo "GoDebut deploy finalizado !"
echo $HASH
echo " "

