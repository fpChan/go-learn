package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	// 你需要完成这个函数。你可与获取到来自map任务生产的中间数据，通过reduceName获取到文件名。
	// 记住你应该编码了值到中间文件,所以你需要解码它们。如果你选择了使用JSON,你通过创建decoder读取到多个
	// 解码之后的值，直接调用Decode直到返回错误。
	//
	// 你应该将reduce输出以JSON编码的方式保存到文件，文件名通过mergeName获取。我们建议你在这里使用JSON,

	// key是中间文件里面键值，value是字符串,这个map用于存储相同键值元素的合并

	// Reduce的过程如下：
	//  S1: 获取到Map产生的文件并打开(reduceName获取文件名)
	// 　S2：获取中间文件的数据(对多个map产生的文件更加值合并)
	// 　S3：打开文件（mergeName获取文件名），将用于存储Reduce任务的结果
	// 　S4：合并结果之后(S2)，进行reduceF操作, work count的操作将结果累加，也就是word出现在这个文件中出现的次数

	var keys []string
	var kvs = make(map[string][]string) //存储所有的key-value对
	for i := 0; i < nMap; i++ {
		fn := reduceName(jobName, i, reduceTask)
		tmpFile, err := os.Open(fn)
		if err != nil {
			log.Printf("open immediate file %s failed", fn)
			continue
		}
		var kv KeyValue

		decoder := json.NewDecoder(tmpFile)
		err = decoder.Decode(&kv)
		for err == nil {
			if _, ok := kvs[kv.Key]; !ok {
				keys = append(keys, kv.Key)
			}
			kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
			err = decoder.Decode(&kv)
		}
	}
	sort.Strings(keys)
	out, err := os.Create(outFile)
	if err != nil {
		log.Printf("create output file %s failed", outFile)
		return
	}
	encoder := json.NewEncoder(out)
	for _, key := range keys {
		if err = encoder.Encode(KeyValue{key, reduceF(key, kvs[key])}); err != nil {
			log.Printf("write [key: %s] to file %s failed", key, outFile)
		}
	}
	_ = out.Close()

}
