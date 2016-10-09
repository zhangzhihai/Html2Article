<?php
/**
 * 采集正文
 */

class htmlarticle{

	public $preTextLen = 0; // 记录上一次统计的字符数量
	public $startPos = - 1; // 记录文章正文的起始位置
	public $_depth = 6; // 按行分析的深度，默认为6 0-6行满足limitcount
	public $_limitCount = 180; // 字符限定数，当分析的文本数量达到限定数则认为进入正文内容
	public $__headEmptyLines = 2; // 确定文章正文头部时，向上查找，连续的空行到达_headEmptyLines，则停止查找
	public $_endLimitCharCount = 20; // 用于确定文章结束的字符数
	public $_appendMode = false;//
	public $sb = array ();//剔除标签分割内容
	public $orgSb = array ();//未出剔除标签分割内容
	public $_data = null;//正文内容
	
	/**
	 * 从body标签文本中分析正文内容
	 */
	function GetContent(){
		$content = $this->_data;
		if (substr_count ( $content, "\n" ) < 10) {
			$content = str_replace ( ">", ">\n", $content );
		}
		
		$regex = '/<body(.*?)<\/body>/is';
		preg_match_all ( $regex, $content, $matches );
		//var_dump($content);
		
		$body = $matches [0] [0];
		// 过滤不相干标签
		$body = $this->striptags ( $body );
		// 保存原始内容，按行存储
		$orgLines = explode ( "\n", $body );
		
		foreach ( $orgLines as $k => $v ) {
			$lines [$k] = $v ;
		}
		
		//处理分割
		for($i = 0; $i < count ( $lines ) - $this->_depth; $i ++) {
			$len = 0;
			//字符串长度
			for($j = 0; $j < $this->_depth; $j ++) {
				$len += strlen ( $lines [$i + $j] );
			}
			//var_dump($len);
		
			// 还没有找到文章起始位置，需要判断起始位置
			if ($this->startPos == - 1) {
				
				if ($this->preTextLen > $this->_limitCount && $len > 0){
					// 如果上次查找的文本数量超过了限定字数，且当前行数字符数不为0，则认为是开始位置
					// 查找文章起始位置, 如果向上查找，发现2行连续的空行则认为是头部
					$emptyCount=0;
					
					for($j = $i - 1; $j > 0; $j --) {
						if(!isset($lines [$j])){
							continue;
						}
						if ($lines [$j]) {
		
							$emptyCount++;
						} else {
							$emptyCount=0;
						}
						if ($emptyCount == $this->__headEmptyLines) {
							$this->startPos = $j + $this->__headEmptyLines;
							break;
						}
					}
					
					// 如果没有定位到文章头，则以当前查找位置作为文章头
					if ($this->startPos == - 1) {
						$this->startPos = $i;
					}
					// 填充发现的文章起始部分
					for($j = $this->startPos; $j <= $i; $j ++) {
						array_push ( $this->sb, $lines [$i] );
						array_push ( $this->orgSb, $orgLines [$i] );
					}
				}
		
			} else {
				
				// 当前长度为0，且上一个长度也为0，则认为已经结束
				if ($len <= $this->_endLimitCharCount && $this->preTextLen < $this->_endLimitCharCount)
				{
					//追加模式
					if (! $this->_appendMode) {
						break;
					}
					$startPos = - 1;
				}
				array_push ( $this->sb, $lines [$i] );
				array_push ( $this->orgSb, $orgLines [$i] );
			}
			
			$this->preTextLen = $len;
		}
		
		
	}
	
	/**
	 * 获取文章发布日期
	 */
	function GetPublishDate(){
		
	}
	/**
	 * 获取title
	 */
	function GetTitle(){
		$content = $this->_data;
		$regex = '/<title>(.*?)<\/title>/is';
		preg_match_all ( $regex, $content, $matches );
		return $matches [1][0];
	}
	
	/**
	 * 从给定的Html原始文本中获取正文信息
	 */
	function GetArticle($content=null){
		$this->_data = $content;
		$this->GetContent();
		
	}
	
	/**
	 * 过滤所有标签保存纯文本
	 */
	function striptags($content) {
		$search = array (
				"'<(\/?)div[^>]*>'si",
				"'<style[^>]*?>.*?</style>'si",
				"'<!--.*?-->'si",
				"'<script[^>]*?>.*?</script>'si",
				"'</?iframe[^>]*?>'si",
				"'<\?[^>]*?\?>'si",
				"'&lt;/?iframe[^>]*?&gt;'si"
		);
		$strip = preg_replace ( $search, "", $content );
		$search = "'<[^>]*>'si"; /* strip html tag */
		return $strip = preg_replace ( $search, "", $strip );
	}
}

set_time_limit(0);

$content = file_get_contents ( '' );


$htmlarticle = new htmlarticle();
$htmlarticle->GetArticle($content);

//echo join("\n",$htmlarticle->orgSb);
