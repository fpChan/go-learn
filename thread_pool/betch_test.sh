export GO111MODULE=on
go mod tidy
Analysis=analysis

if [ ! -d ${Analysis} ]; then
    mkdir $Analysis
fi

cd src

validCPUFile=../${Analysis}/validcpu.prof
validMEMFile=../${Analysis}/validmem.pprof

InvalidCPUFile=../${Analysis}/invalidcpu.prof
InvalidMEMFFile=../${Analysis}/invalidmem.pprof

function deleteFile() {
    if  [  -f  "${1}" ]; then
      echo "$1"
      all=*
      rm  -f  ${1}${all}
    fi
}


function clear() {
  files=("${validCPUFile}" "${validMEMFile}" "${InvalidCPUFile}" "${InvalidMEMFFile}" "./poolhash" "./nativehash")
  for file in "${files[@]}"
  do
    deleteFile "$file"
  done
}

function performance_analysis () {

    go test -bench='BenchmarkPoolGoroutineAtomicAdd' -benchmem  -cpuprofile=${validCPUFile}  --memprofile=${validMEMFile}
    go test -bench='BenchmarkPoolGoruntinueFileIO' -benchmem  -cpuprofile=${InvalidCPUFile}  --memprofile=${InvalidMEMFFile}
    echo svg | go tool pprof ${validCPUFile} && mv profile001.svg ../${Analysis}/validCPUFile.svg
    echo svg | go tool pprof ${validMEMFile} && mv profile001.svg ../${Analysis}/validMEMFile.svg
    echo svg | go tool pprof ${InvalidCPUFile} && mv profile001.svg ../${Analysis}/InvalidCPUFile.svg
    echo svg | go tool pprof ${InvalidMEMFFile} && mv profile001.svg ../${Analysis}/InvalidMEMFFile.svg
}

function exec_test() {
    echo "The effect of pool test is obvious:  Atomic Add OP"
    go test -bench='AtomicAdd' -benchmem

    echo "The effect of pool test is obvious:  File IO"
    go test -bench='FileIO' -benchmem
}


clear
#performance_analysis
exec_test