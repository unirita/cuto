#/bin/sh

go get github.com/coopernurse/gorp
if [ "$?" -ne 0 ] ; then
  echo go get error.
  exit 1
fi

go get github.com/BurntSushi/toml
if [ "$?" -ne 0 ] ; then
  echo go get error.
  exit 1
fi

go get github.com/cihub/seelog
if [ "$?" -ne 0 ] ; then
  echo go get error.
  exit 1
fi

go get github.com/mattn/go-sqlite3
if [ "$?" -ne 0 ] ; then
  echo go get error.
  exit 1
fi

go get golang.org/x/tools/cmd/cover
if [ "$?" -ne 0 ] ; then
  echo go get error.
  exit 1
fi

exit 0

