
#!/bin/bash

./bin/uhppote-cli set-door-control 423187757 1 'normally open'  
./bin/uhppote-cli set-door-control 423187757 2 'normally open'  
./bin/uhppote-cli set-door-control 423187757 3 'normally open'  
./bin/uhppote-cli set-door-control 423187757 4 'normally open'  

count=1
while [ $count -le 1 ]
do
    echo "${count}"
	./bin/uhppote-cli open 423187757 1                            
	./bin/uhppote-cli open 423187757 2                            
	./bin/uhppote-cli open 423187757 3                            
	./bin/uhppote-cli open 423187757 4                            
    sleep 1
	((count++))
done

./bin/uhppote-cli set-door-control 423187757 1 'controlled'  
./bin/uhppote-cli set-door-control 423187757 2 'controlled'  
./bin/uhppote-cli set-door-control 423187757 3 'controlled'  
./bin/uhppote-cli set-door-control 423187757 4 'controlled'  

./bin/uhppote-cli open 423187757 1                            
./bin/uhppote-cli open 423187757 2                            
./bin/uhppote-cli open 423187757 3                            
./bin/uhppote-cli open 423187757 4                            

events=$(./bin/uhppote-cli get-events 423187757)
first=$(./bin/uhppote-cli get-event  423187757 first)
last=$(./bin/uhppote-cli get-event  423187757 last)

echo "$(date)"           >> events.log
echo "EVENTS: ${events}" >> events.log
echo "FIRST:  ${first}"  >> events.log
echo "LAST:   ${last}"   >> events.log
echo ""                  >> events.log

tail -n 12 events.log

say 'woooot'
