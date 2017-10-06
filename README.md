# comm
Message communication SDK base on golang

## 性能测试
### 测试环境
CPU: Intel(R)Core(TM)m3-6y30 CPU@ 0.9GHz 1.51GHz<br>
RAM: 4GB<br>

### 测试数据

<table>
    <tr>
        <th> 消息长度  </th>
        <th> 吞吐量<br>（单位：ktps） </th>
        <th> 流量<br>（MB/s） </th>
    </tr>
    <tr>
        <th>8</th>
        <th>1069</th>
        <th>8.532</th>
    </tr>
    <tr>
        <th>16</th>
        <th>921</th>
        <th>14.777</th>
    </tr>
    <tr>
        <th>32</th>
        <th>897</th>
        <th>28.683</th>
    </tr>
    <tr>
        <th>64</th>
        <th>836</th>
        <th>53.372</th>
    </tr>
    <tr>
        <th>128</th>
        <th>701</th>
        <th>89.122</th>
    </tr>
	<tr>
        <th>256</th>
        <th>500</th>
        <th>126.242</th>
    </tr>
	<tr>
        <th>512</th>
        <th>247</th>
        <th>124.85</th>
    </tr>
	<tr>
        <th>1024</th>
        <th>116</th>
        <th>116.882</th>
    </tr>
	<tr>
        <th>2048</th>
        <th>61</th>
        <th>122.215</th>
    </tr>
	<tr>
        <th>4096</th>
        <th>40</th>
        <th>161.793</th>
    </tr>
	<tr>
        <th>8192</th>
        <th>22</th>
        <th>177.25</th>
    </tr>
	<tr>
        <th>16384</th>
        <th>13</th>
        <th>223.516</th>
    </tr>
	<tr>
        <th>32768</th>
        <th>9</th>
        <th>294.281</th>
    </tr>
</table>
