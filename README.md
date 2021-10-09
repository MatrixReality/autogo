Autobot-go WIP
Motores funcionais

Commit inicial

Em uma plataforma raspberry que siga a pinagem do esquema autobot, instalar o pi-blaster
https://github.com/sarfata/pi-blaster

Ativar gpio com 'pins enable' na pasta de instalação do pi-blaster:
'sudo pi-blaster --gpio 17,18,16,19,13,20,12,21,22,5,23,4,24,25,26,27'

Gerando Binario com raspberry como device alvo:
'GOARM=6 GOARCH=arm GOOS=linux go build main.go'

Rode o Binario, seja feliz

Referências:
https://gobot.io/documentation/platforms/raspi/
https://gobot.io/documentation/examples/firmata_motor/
https://pkg.go.dev/github.com/heupel/gobot/platforms/gpio#section-readme

-Próximas etapas (deixar identico ao ultimo Master da versão Python):
  - Controle via teclado
  - Modulo LCD funcional
  - Refatoração na estrutura do código
  - Comunicação com arduino (Sonar set)
  - sh e makefile para automatizar dependencias em instalação nova
  - sh update de goversion no raspbian