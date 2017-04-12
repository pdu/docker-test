for ((p=10000;p<10100;p++))
do
    nohup ./server -port=":$p" &
done
