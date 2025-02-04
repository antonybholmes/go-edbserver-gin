pwd=`pwd`
for d in `find . -maxdepth 1 -mindepth 1 -type d`
do
	echo ${d}
	cd ${d}
	go get -u
	cd ${pwd}
done
