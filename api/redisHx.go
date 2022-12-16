package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func Redistd(c *gin.Context) {
	ctx := context.Background()

	setGetExample(redisClient, ctx)
	stringIntExample(redisClient, ctx)
	listExample(redisClient, ctx)
	setExample(redisClient, ctx)
	zsetExample(redisClient, ctx)
	hashExample(redisClient, ctx)
	hyperLogLogExample(redisClient, ctx)
	ExampleClient_CMD(redisClient, ctx)
	TxPipelineExample(redisClient, ctx)
	ScriptExample(redisClient, ctx)

}

func setGetExample(rdb *redis.Client, ctx context.Context) {
	// 1.Set 设置 key 如果设置为-1则表示永不过期
	err := rdb.Set(ctx, "score", 100, 60*time.Second).Err()
	if err != nil {
		basicLog.Errorf("set score failed, err:%v\n", err)
		panic(err)
	}

	// 2.Get 获取已存在的Key其存储的值
	val1, err := rdb.Get(ctx, "score").Result() // 获取其值
	if err != nil {
		basicLog.Errorf("get score failed, err:%v\n", err)
		panic(err)
	}
	basicLog.Debugf("val1 -> score ：%v\n", val1)

	// Get 获取一个不存在的值返回redis.Nil 则说明不存在
	val2, err := rdb.Get(ctx, "name").Result()
	if err == redis.Nil {
		basicLog.Errorf("[ERROR] - Key [name] not exist")
	} else if err != nil {
		basicLog.Errorf("get name failed, err:%v\n", err)
		panic(err)
	}
	// Exists() 方法用于检测某个key是否存在
	n, _ := rdb.Exists(ctx, "name").Result()
	if n > 0 {
		basicLog.Debugln("name key 存在!")
	} else {
		basicLog.Error("name key 不存在!")
		rdb.Set(ctx, "name", "weiyi", 60*time.Second)
	}
	val2, _ = rdb.Get(ctx, "name").Result()
	basicLog.Debugln("val2 -> name : ", val2)

	// 3.SetNX 当不存在key时将进行设置该可以并设置其过期时间
	val3, err := rdb.SetNX(ctx, "username", "weiyigeek", 0).Result()
	if err != nil {
		basicLog.Errorf("set username failed, err:%v\n", err)
		panic(err)
	}
	basicLog.Debugln("val3 -> username: %v\n", val3)

	// 4.Keys() 根据正则获取keys, DBSize() 查看当前数据库key的数量.
	keys, _ := rdb.Keys(ctx, "*").Result()
	num, err := rdb.DBSize(ctx).Result()
	if err != nil {
		panic(err)
	}
	basicLog.Debugln("All Keys : %v, Keys number : %v \n", keys, num)

	// 根据前缀获取Key
	vals, _ := rdb.Keys(ctx, "user*").Result()
	basicLog.Debugln(vals)

	// 5.Type() 方法用户获取一个key对应值的类型
	vType, err := rdb.Type(ctx, "username").Result()
	if err != nil {
		panic(err)
	}
	basicLog.Debugln("username key type : %v\n", vType)

	// 6.Expire()方法是设置某个时间段(time.Duration)后过期，ExpireAt()方法是在某个时间点(time.Time)过期失效.
	val4, _ := rdb.Expire(ctx, "name", time.Minute*2).Result()
	if val4 {
		basicLog.Debugln("name 过期时间设置成功", val4)
	} else {
		basicLog.Error("name 过期时间设置失败", val4)
	}
	val5, _ := rdb.ExpireAt(ctx, "username", time.Now().Add(time.Minute*2)).Result()
	if val5 {
		basicLog.Debugln("username 过期时间设置成功", val5)
	} else {
		basicLog.Error("username 过期时间设置失败", val5)
	}

	// 7.TTL()与PTTL()方法可以获取某个键的剩余有效期
	userTTL, _ := rdb.TTL(ctx, "user").Result() // 获取其key的过期时间
	usernameTTL, _ := rdb.PTTL(ctx, "username").Result()
	basicLog.Debugf("user TTL : %v, username TTL : %v\n", userTTL, usernameTTL)

	// 8.Del():删除缓存项与FlushDB():清空当前数据
	// 当通配符匹配的key的数量不多时，可以使用Keys()得到所有的key在使用Del命令删除。
	num, err = rdb.Del(ctx, "user", "username").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Del() : ", num)
	// 如果key的数量非常多的时候，我们可以搭配使用Scan命令和Del命令完成删除。
	iter := rdb.Scan(ctx, 0, "user*", 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	// 9.清空当前数据库，因为连接的是索引为0的数据库，所以清空的就是0号数据库
	flag, err := rdb.FlushDB(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("FlushDB() : ", flag)

}

func stringIntExample(rdb *redis.Client, ctx context.Context) {
	// 设置字符串类型的key
	err := rdb.Set(ctx, "hello", "Hello World!", 0).Err()
	if err != nil {
		panic(err)
	}
	// GetRange ：字符串截取
	// 注：即使key不存在，调用GetRange()也不会报错，只是返回的截取结果是空"",可以使用fmt.Printf("%q\n", val)来打印测试
	val1, _ := rdb.GetRange(ctx, "hello", 1, 4).Result()
	basicLog.Debugf("key: hello, value: %v\n", val1) //截取到的内容为: ello

	// Append()表示往字符串后面追加元素，返回值是字符串的总长度
	length1, _ := rdb.Append(ctx, "hello", " Go Programer").Result()
	val2, _ := rdb.Get(ctx, "hello").Result()
	basicLog.Debugf("当前缓存key的长度为: %v，值: %v \n", length1, val2)

	// 设置整形的key
	err = rdb.SetNX(ctx, "number", 1, 0).Err()
	if err != nil {
		panic(err)
	}
	// Incr()、IncrBy()都是操作数字，对数字进行增加的操作
	// Decr()、DecrBy()方法是对数字进行减的操作，和Incr正好相反
	// incr是执行原子加1操作
	val3, _ := rdb.Incr(ctx, "number").Result()
	basicLog.Debugf("Incr -> key当前的值为: %v\n", val3) // 2
	// incrBy是增加指定的数
	val4, _ := rdb.IncrBy(ctx, "number", 6).Result()
	basicLog.Debugf("IncrBy -> key当前的值为: %v\n", val4) // 8

	// StrLen 也可以返回缓存key的长度
	length2, _ := rdb.StrLen(ctx, "number").Result()
	basicLog.Debugf("number 值长度: %v\n", length2)
}

func listExample(rdb *redis.Client, ctx context.Context) {
	// 插入指定值到list列表中，返回值是当前列表元素的数量
	// 使用LPush()方法将数据从左侧压入链表（后进先出）,也可以从右侧压如链表对应的方法是RPush()
	count, _ := rdb.LPush(ctx, "list", 1, 2, 3).Result()
	fmt.Println("插入到list集合中元素的数量: ", count)

	// LInsert() 在某个位置插入新元素
	// 在名为key的缓存项值为2的元素前面插入一个值，值为123 ， 注意只会执行一次
	_ = rdb.LInsert(ctx, "list", "before", "2", 123).Err()
	// 在名为key的缓存项值为2的元素后面插入一个值，值为321
	_ = rdb.LInsert(ctx, "list", "after", "2", 321).Err()

	// LSet() 设置某个元素的值
	//下标是从0开始的
	val1, _ := rdb.LSet(ctx, "list", 2, 256).Result()
	fmt.Println("是否成功将下标为2的元素值改成256: ", val1)

	// LLen() 获取链表元素个数
	length, _ := rdb.LLen(ctx, "list").Result()
	fmt.Printf("当前链表的长度为: %v\n", length)

	// LIndex() 获取链表下标对应的元素
	val2, _ := rdb.LIndex(ctx, "list", 2).Result()
	fmt.Printf("下标为2的值为: %v\n", val2)

	// 从链表左侧弹出数据
	val3, _ := rdb.LPop(ctx, "list").Result()
	fmt.Printf("弹出下标为0的值为: %v\n", val3)

	// LRem() 根据值移除元素 lrem key count value
	n, _ := rdb.LRem(ctx, "list", 2, "256").Result()
	fmt.Printf("移除了: %v 个\n", n)
}

func setExample(rdb *redis.Client, ctx context.Context) {
	// 集合元素缓存设置
	keyname := "Program"
	mem := []string{"C", "Golang", "C++", "C#", "Java", "Delphi", "Python", "Golang"}
	// //由于Golang已经被添加到Program集合中，所以重复添加时无效的
	for _, v := range mem {
		rdb.SAdd(ctx, keyname, v)
	}

	// SCard() 获取集合元素个数
	total, _ := rdb.SCard(ctx, keyname).Result()
	fmt.Println("golang集合成员个数: ", total)

	// SPop() 随机获取一个元素 （无序性，是随机的）
	val1, _ := rdb.SPop(ctx, keyname).Result()
	// SPopN()  随机获取多个元素.
	val2, _ := rdb.SPopN(ctx, keyname, 2).Result()

	// SSMembers() 获取所有成员
	val3, _ := rdb.SMembers(ctx, keyname).Result()
	fmt.Printf("随机获取一个元素: %v , 随机获取多个元素: %v \n所有成员: %v\n", val1, val2, val3)

	// SIsMember() 判断元素是否在集合中
	exists, _ := rdb.SIsMember(ctx, keyname, "golang").Result()
	if exists {
		fmt.Println("golang 存在 Program 集合中.") // 注意:我们存入的是Golang而非golang
	} else {
		fmt.Println("golang 不存在 Program 集合中.")
	}

	// SUnion():并集, SDiff():差集, SInter():交集
	rdb.SAdd(ctx, "setA", "a", "b", "c", "d")
	rdb.SAdd(ctx, "setB", "a", "d", "e", "f")

	//并集
	union, _ := rdb.SUnion(ctx, "setA", "setB").Result()
	fmt.Println("并集", union)

	//差集
	diff, _ := rdb.SDiff(ctx, "setA", "setB").Result()
	fmt.Println("差集", diff)

	//交集
	inter, _ := rdb.SInter(ctx, "setA", "setB").Result()
	fmt.Println("交集", inter)

	// 删除集合中指定元素(返回成功)
	n, _ := rdb.SRem(ctx, "setB", "a", "f").Result()
	fmt.Println("已成功删除元素的个数: ", n)
}

func zsetExample(rdb *redis.Client, ctx context.Context) {
	// 有序集合成员与分数设置
	// zSet类型需要使用特定的类型值*redis.Z，以便作为排序使用
	lang := []*redis.Z{
		&redis.Z{Score: 90.0, Member: "Golang"},
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C/C++"},
	}
	//插入ZSet类型
	num, err := rdb.ZAdd(ctx, "language_rank", lang...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Printf("zadd %d succ.\n", num)

	// 将ZSet中的某一个元素顺序值增加: 把Golang的分数加10
	newScore, err := rdb.ZIncrBy(ctx, "language_rank", 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 根据分数排名取出元素:取分数最高的3个
	ret, err := rdb.ZRevRangeWithScores(ctx, "language_rank", 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return
	}
	fmt.Printf("zsetKey前3名热度的是: %v\n,Top 3 的 Memeber 与 Score 是:\n", ret)
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// ZRangeByScore()、ZRevRangeByScore():获取score过滤后排序的数据段
	// 此处表示取95~100分的
	op := redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, "language_rank", &op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	// 输出全部成员及其score分数
	fmt.Println("language_rank 键存储的全部元素:")
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

func hashExample(rdb *redis.Client, ctx context.Context) {
	// (1) HSet() 设置字段和值
	rdb.HSet(ctx, "huser", "key1", "value1", "key2", "value2")
	rdb.HSet(ctx, "huser", []string{"key3", "value3", "key4", "value4"})
	rdb.HSet(ctx, "huser", map[string]interface{}{"key5": "value5", "key6": "value6"})

	// (2) HMset():批量设置
	rdb.HMSet(ctx, "hmuser", map[string]interface{}{"name": "WeiyiGeek", "age": 88, "address": "重庆"})

	// (3) HGet() 获取某个元素
	address, _ := rdb.HGet(ctx, "hmuser", "address").Result()
	fmt.Println("hmuser.address -> ", address)

	// (4) HGetAll() 获取全部元素
	hmuser, _ := rdb.HGetAll(ctx, "hmuser").Result()
	fmt.Println("hmuser :=> ", hmuser)

	// (5) HExists 判断元素是否存在
	flag, _ := rdb.HExists(ctx, "hmuser", "address").Result()
	fmt.Println("address 是否存在 hmuser 中: ", flag)

	// (6) HLen() 获取长度
	length, _ := rdb.HLen(ctx, "hmuser").Result()
	fmt.Println("hmuser hash 键长度: ", length)

	// (7) HDel() 支持一次删除多个元素
	count, _ := rdb.HDel(ctx, "huser", "key3", "key4").Result()
	fmt.Println("删除元素的个数: ", count)
}

func hyperLogLogExample(rdb *redis.Client, ctx context.Context) {
	//基数统计 HyperLogLog 类型操作
	//描述: 用来做基数统计的算法，HyperLogLog 的优点是，在输入元素的数量或者体积非常非常大时，计算基数所需的空间总是固定 的、并且是很小的。
	//Tips: 每个 HyperLogLog 键只需要花费 12 KB 内存，就可以计算接近 2^64 个不同元素的基数.

	log.Println("Start ExampleClient_HyperLogLog")
	defer log.Println("End ExampleClient_HyperLogLog")
	//  设置 HyperLogLog 类型的键  pf_test_1
	for i := 0; i < 5; i++ {
		rdb.PFAdd(ctx, "pf_test_1", fmt.Sprintf("pf1key%d", i))
	}
	ret, err := rdb.PFCount(ctx, "pf_test_1").Result()
	log.Println(ret, err)

	//  设置 HyperLogLog 类型的键  pf_test_2
	for i := 0; i < 10; i++ {
		rdb.PFAdd(ctx, "pf_test_2", fmt.Sprintf("pf2key%d", i))
	}
	ret, err = rdb.PFCount(ctx, "pf_test_2").Result()
	log.Println(ret, err)

	//  合并两个 HyperLogLog 类型的键  pf_test_1 + pf_test_1
	rdb.PFMerge(ctx, "pf_test", "pf_test_2", "pf_test_1")
	ret, err = rdb.PFCount(ctx, "pf_test").Result()
	log.Println(ret, err)
}

func ExampleClient_CMD(rdb *redis.Client, ctx context.Context) {
	//自定义redis指令操作
	//描述: 我们可以采用go-redis提供的Do方法，可以让我们直接执行redis-cli中执行的相关指令, 可以极大的便于使用者上手。

	log.Println("Start ExampleClient_CMD")
	defer log.Println("End ExampleClient_CMD")

	// 1.执行redis指令 Set 设置缓存
	v := rdb.Do(ctx, "set", "NewStringCmd", "redis-cli").String()
	log.Println(">", v)

	// 2.执行redis指令 Get 设置缓存
	v = rdb.Do(ctx, "get", "NewStringCmd").String()
	log.Println("Method1 >", v)

	// 3.匿名方式执行自定义redis命令
	// Set
	Set := func(client *redis.Client, ctx context.Context, key, value string) *redis.StringCmd {
		cmd := redis.NewStringCmd(ctx, "set", key, value) // 关键点
		client.Process(ctx, cmd)
		return cmd
	}
	v, _ = Set(rdb, ctx, "NewCmd", "go-redis").Result()
	log.Println("> set NewCmd go-redis:", v)

	// Get
	Get := func(client *redis.Client, ctx context.Context, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd(ctx, "get", key) // 关键点
		client.Process(ctx, cmd)
		return cmd
	}
	v, _ = Get(rdb, ctx, "NewCmd").Result()
	log.Println("Method2 > get NewCmd:", v)

	// 4.执行redis指令 hset 设置哈希缓存 (实践以下方式不行)
	// kv := map[string]interface{}{"key5": "value5", "key6": "value6"}
	// v, _ = rdb.Do(ctx, "hmset", "NewHashCmd", kv)
	// log.Println("> ", v)
}

func TxPipelineExample(rdb *redis.Client, ctx context.Context) {
	//所以在某些场景下，当我们有多条命令要执行时，就可以考虑使用pipeline来优化redis缓冲效率。

	//MULTI/EXEC 事务处理操作
	//描述: Redis是单线程的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，例如在它们之间交替执行。但是Multi/exec
	//能够确保在其两个语句之间的命令之间没有其他客户端正在执行命令。
	//
	//在这种场景我们需要使用TxPipeline, 它总体上类似于上面的Pipeline, 但是它内部会使用MULTI/EXEC
	//包裹排队的命令。例如：

	//pipe := rdb.TxPipeline()
	//incr := pipe.Incr("tx_pipeline_counter")
	//pipe.Expire("tx_pipeline_counter", time.Hour)
	//_, err := pipe.Exec()
	//fmt.Println(incr.Val(), err)

	// # 上面代码相当于在一个RTT下执行了下面的redis命令：
	//MULTI
	//INCR pipeline_counter
	//EXPIRE pipeline_counts 3600
	//EXEC

	// 开pipeline与事务
	pipe := rdb.TxPipeline()
	// 设置TxPipeline键缓存
	v, _ := rdb.Do(ctx, "set", "TxPipeline", 1023.0).Result()
	log.Println(v)
	// 自增+1.0
	incr := pipe.IncrByFloat(ctx, "TxPipeline", 1026.0)
	log.Println(incr) // 未提交时  incr.Val() 值 为 0
	// 设置键过期时间
	pipe.Expire(ctx, "TxPipeline", time.Hour)
	// 提交事务
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("执行失败, 进行回滚操作!")
		return
	}
	fmt.Println("事务执行成功,已提交!")
	log.Println("TxPipeline :", incr.Val()) // 提交后值 为 2049
}

func ScriptExample(rdb *redis.Client, ctx context.Context) {
	//描述: 从 Redis 2.6.0 版本开始的，使用内置的 Lua 解释器，可以对 Lua 脚本进行求值, 所以我们可直接在redis客户端中执行一些脚本。
	//
	//redis Eval 命令基本语法如下：EVAL script numkeys key [key ...] arg [arg ...]
	//
	//script: 参数是一段 Lua 5.1 脚本程序。脚本不必(也不应该)定义为一个 Lua 函数。
	//numkeys: 用于指定键名参数的个数。
	//key [key ...]: 从 EVAL 的第三个参数开始算起，表示在脚本中所用到的那些 Redis 键(key)，这些键名参数可以在 Lua 中通过全局变量 KEYS 数组，用 1 为基址的形式访问( KEYS[1] ， KEYS[2] ，以此类推)
	//。
	//arg [arg ...]: 附加参数，在 Lua 中通过全局变量 ARGV 数组访问，访问的形式和 KEYS 变量类似( ARGV[1] 、 ARGV[2] ，诸如此类)
	//。
	//redis.call()
	//与 redis.pcall()
	//唯一的区别是当redis命令执行结果返回错误时 redis.call() 将返回给调用者一个错误，而redis.pcall()会将捕获的错误以Lua表的形式返回

	// Lua脚本定义1. 传递key输出指定格式的结果
	EchoKey := redis.NewScript(`
		if redis.call("GET", KEYS[1]) ~= false then
			return {KEYS[1],"==>",redis.call("get", KEYS[1])}
		end
		return false
	`)

	err := rdb.Set(ctx, "xx_name", "WeiyiGeek", 0).Err()
	if err != nil {
		panic(err)
	}
	val1, err := EchoKey.Run(ctx, rdb, []string{"xx_name"}).Result()
	log.Println(val1, err)

	// Lua脚本定义2. 传递key与step使得，key值等于`键值+step`
	IncrByXX := redis.NewScript(`
		if redis.call("GET", KEYS[1]) ~= false then
			return redis.call("INCRBY", KEYS[1], ARGV[1])
		end
		return false
	`)

	// 判断键是否存在，存在就删除该键
	exist, err := rdb.Exists(ctx, "xx_counter").Result()
	if exist > 0 {
		res, err := rdb.Del(ctx, "xx_counter").Result()
		log.Printf("is Exists?: %v, del xx_counter: %v, err: %v \n", exist, res, err)
	}

	// 首次调用
	val2, err := IncrByXX.Run(ctx, rdb, []string{"xx_counter"}, 2).Result()
	log.Println("首次调用 IncrByXX.Run ->", val2, err)

	// 写入 xx_counter 键
	err = rdb.Set(ctx, "xx_counter", 40, 0).Err()
	if err != nil {
		panic(err)
	}
	// 二次调用
	val3, err := IncrByXX.Run(ctx, rdb, []string{"xx_counter"}, 2).Result()
	log.Println("二次调用 IncrByXX.Run ->", val3, err)
}
