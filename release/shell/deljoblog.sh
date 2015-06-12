#/bin/sh

# ################
GENERAL_NUM=30
# ################

if [ -z "$CUTOROOT" ] ; then
  echo Undefined \$CUTOROOT
  exit $RC
fi

files="$CUTOROOT/joblog/*"
dirary=()
for filepath in $files; do
  if [ -d $filepath ] ; then
    dirary+=("$filepath")
  fi  
done    

delnum=`expr ${#dirary[*]} - $GENERAL_NUM`
if [ $delnum -lt 1 ] ; then
  echo Do not delete.
  exit 0
fi

count=1
for dir in ${dirary[@]}; do
  rm -r $dir
  echo  deleted directory [ $dir ]
  count=`expr $count + 1`
  if [ $count -gt $delnum ] ; then
    break
  fi  
done

exit 0
