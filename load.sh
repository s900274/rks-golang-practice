#!/bin/bash
# progdir=$(dirname $0);
CONFIGNAME=$1
cd `dirname $0` || exit
#ulimit -c unlimited
mkdir -p /apps/engine/magneto/status/magneto
mkdir -p /apps/engine/magneto/log

start(){
	stop
	sleep 1
	setsid /apps/engine/magneto/bin/supervise.magneto -u /apps/engine/magneto/status/magneto env GOTRACEBACK=crash /apps/engine/magneto/bin/magneto -config /apps/engine/magneto/config/magneto.${CONFIGNAME}.toml
}

stop(){
	killall -9 supervise.magneto
	killall -9 magneto
}

restart(){
	killall -9 magneto
}


case C"$2" in
	Cstart)
		start
		echo "start Done!"
		;;
	Cstop)
		stop
		echo "stop Done!"
		;;
	Crestart)
		restart
		echo "restart Done!"
		;;
	C*)
		echo "Usage: $0 {start|stop|restart}"
		;;
esac
