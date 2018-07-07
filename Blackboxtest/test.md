## 测试文档

- ### **测试总结**<br>
编写了一个黑盒测试脚本，可以做到不依赖数据库之前有的数据即以往的资源按照一定的顺序自动化完成对DeliciousFoodEasyOrder中的全部api的测试。<br>

- ### **测试逻辑设计**<br>
如果想重复的且不依赖以往的资源并只输入一行指令就可以一次性将api文档中的所有api都测试一遍，那么访问api的顺序就一定不是按照api文档的顺序从头到尾访问一遍列出的api。
逻辑设计为先访问Create a merchant这个api，成功以后可以根据注册的merchant来完成认证即访问Password Grant这个api。完成对这两个的api的测试以后就可以得到token了。有了token以后就能成功访问需要携带token才能成功访问的api了。之后是将剩下的关于merchant的api都测试一遍。

完成对merchant相关的api的测试以后，就可以逐个完成customer相关的api,seat相关的api和food
相关的api的测试，最后再完成order部分的api测试，因为order部分的很多请求都需要传入customer _id,food_id,seat_id这种data。

在所有的测试中都是create操作最先执行，delete操作最后执行(有的部分没有delete的api就不执行)。因为该测试脚本要做到不依赖以往的资源完成对全部api的访问。所以要自行创建资源然后进行crud中的cru操作，最后进行crud中的d操作。

- ### **其他说明**<br>
如果中途有哪个api的访问结果出错了，即服务器不返回报文或者返回的报文不是预期的报文，那么黑盒测试程序会终止(exit)，并且会在对应的地方报错，直接将错误输出到控制台上，便于debug。每成功访问一个api，黑盒测试程序也会在控制台上输出成功的信息，如果所有api都没有错误，程序会输出对应数目的对应api名字的成功信息一直到程序自动结束（而不是中途exit）。
