# utils

工具包

### bitmap

位图。用一个 bit 位来标记某个元素对应的 Value，而 Key 即是该元素。 由于采用了 Bit 为单位来存储数据，
因此在内存占用方面，可以大大节省。常用于表示一个数据是否出现过，0为没有出现过，1表示出现过。

```
bm := New(8)
t.Log(bm.Cap(), bm.Len()) // 8 0
t.Log(bm.Set(0), bm.Set(6), bm.Set(6), bm.Len()) // true true false 2
t.Log(bm.Set(8), bm.Set(9), bm.Len()) // false false 2, 超过范围
t.Log(bm.Set(2), bm.Clear(2), bm.Clear(3), bm.Len()) // true true false 2
t.Log(bm.String())  // 10000010
```



###  buffer 

byte环形缓存区，减少字节拷贝次数。

###  heap

堆，使用 container/heap 的封装

###  log 日志

* 支持时间分割，每天切分日志文件。 
* 支持文件大小分割，达到日志存储上限，切分日志。
* 支持日志等级划分。

###  ordermap 
    
有序的map

###  pipeline 

流水线 ， step

###  queue

channelQueue, blockQueue, priorityQueue

###  timer

* heapTimer 小根堆定时器，高精度 timer。调用系统timer，提供统一管理
    
* timingWheel 时间轮定时器，低精度 timer, 最低精度为毫秒。系统ticker驱动。
如果时间轮精度为10ms， 那么他的误差在 （0，10）ms之间。如果一个任务延迟 500ms，那它的执行时间在490～500ms之间。
按平均来讲，出错的概率均等的情况下，那么这个出错可能会延迟或提前最小刻度的一半，在这里就是10ms/2=5ms.
故，时间轮的 tick 单位 在总延迟时间上，应该不足以影响延迟执行函数处理的事务。
