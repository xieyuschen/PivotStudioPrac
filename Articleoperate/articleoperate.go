package Articleoperate

import (
	"PS_m1_ture/Useroperate"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
type Article struct {
	Author string			//用于存储发帖者
	//Whoelse string			//如果做的话，用于评论者评论，若为空则代表为发帖者发的原文，否则代表评论者
	Title string			//用于存储题目
	Summary string			//用于存储内容概要
	Content string			//用于存储文章主体内容
}
//帖子相关操作
//创建帖子，并写入
func WriteArticle(c *gin.Context)  {
	titleinput := c.Query("title")
	authorinput,_ := c.Cookie("Account")
	summaryinput := c.Query("summary")
	contentinput := c.Query("content")
	stmt, err := Useroperate.MysqlDB.Prepare("INSERT INTO articles SET author=?, title=?, summary=?, content=?")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Error in writing an article",
		})
	}else {
		_, err := stmt.Exec(authorinput, titleinput, summaryinput, contentinput)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Failed to store your article.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "You write the article successfully.",
		})
	}
}
//修改帖子content和summary，title(主键)和author不可改
func ReviseArticle(c *gin.Context)  {
	titleinput := c.Query("title")
	summaryinput := c.Query("summary")
	contentinput := c.Query("content")
	result, err := Useroperate.MysqlDB.Exec("UPDATE articles SET summary=?, content=? where title=?", summaryinput, contentinput, titleinput)
	if err != nil{
		fmt.Printf("Revise failed, err:%v\n", err)
		return
	}
	fmt.Println("Revise article successd:", result)
	rowsaffected, err := result.RowsAffected()
	if err != nil{
		fmt.Printf("Get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}
//删除所发的帖子, 需要表中有该title的article且当前账户为创建者账户有权限删除
func DeleteArticle(c *gin.Context)  {
	A := new(Article)
	titleinput := c.Query("title")
	authorinput,_ := c.Cookie("Account")
	row := Useroperate.MysqlDB.QueryRow("select author, title, summary, content from articles where title=?", titleinput)
	if err := row.Scan(&A.Author, &A.Title, &A.Summary, &A.Content); err != nil{
		fmt.Printf("Scan failed. The table does not have this article. err:%v\n", err)
		return
	}
	if A.Author != authorinput{
		fmt.Printf("Scan failed. Your don't have the privilige to delete this article.")
		return
	}
	result, err := Useroperate.MysqlDB.Exec("delete from articles where title=?", titleinput)
	if err != nil{
		fmt.Printf("Delete failed, err:%v\n", err)
		return
	}
	fmt.Println("Delete article successd", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil{
		fmt.Printf("Get RowsAffected failed, err:%v\n",err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}
//查看自己的帖子
func SeeArticles(c *gin.Context)  {
	A := new(Article)
	authorinput,_ := c.Cookie("Account")
	rows, err:= Useroperate.MysqlDB.Query("SELECT author, title, summary, content FROM articles WHERE author=?", authorinput)
	defer func() {
		if rows != nil{
			rows.Close()	//关掉sql连接
		}
	}()
	if err != nil{
		fmt.Printf("Can't see it, err:%v/n", err)
		return
	}
	for rows.Next(){
		err = rows.Scan(&A.Author, &A.Title, &A.Summary, &A.Content)
		if err != nil{
			fmt.Printf("Scan failed, err:%v\n", err)
			return
		}
		fmt.Println("Scan successd:", *A)
	}
}