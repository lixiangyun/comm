# comm
Message communication SDK base on golang

![](https://travis-ci.com/lixiangyun/comm.svg?branch=master)

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
        <th>1189</th>
        <th>9.408</th>
    </tr>
    <tr>
        <th>16</th>
        <th>1042</th>
        <th>16.549</th>
    </tr>
    <tr>
        <th>32</th>
        <th>962</th>
        <th>30.522</th>
    </tr>
    <tr>
        <th>64</th>
        <th>836</th>
        <th>53.008</th>
    </tr>
    <tr>
        <th>128</th>
        <th>737</th>
        <th>93.288</th>
    </tr>
	<tr>
        <th>256</th>
        <th>544</th>
        <th>137.22</th>
    </tr>
	<tr>
        <th>512</th>
        <th>296</th>
        <th>148.663</th>
    </tr>
	<tr>
        <th>1024</th>
        <th>135</th>
        <th>136.198</th>
    </tr>
	<tr>
        <th>2048</th>
        <th>72</th>
        <th>145.229</th>
    </tr>
	<tr>
        <th>4096</th>
        <th>41</th>
        <th>164.289</th>
    </tr>
	<tr>
        <th>8192</th>
        <th>22</th>
        <th>178.633</th>
    </tr>
	<tr>
        <th>16384</th>
        <th>13</th>
        <th>216.891</th>
    </tr>
	<tr>
        <th>32768</th>
        <th>9</th>
        <th>297.125</th>
    </tr>
</table>
