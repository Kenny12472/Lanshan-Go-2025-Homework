package main

import "fmt"

func average(sum int,count int) float64 {
        return float64(sum/count)
       }

func main(){
        var input int
        var add int
        add = 0
        var count1 int
        count1=0
        var a string

    for {
                  fmt.Print("请输入一个整数(输入0结束):")
                  fmt.Scanln(&input)
                        if input !=0 {
                                 add=add+input
                                 count1=count1+1
                         } else { break }
          }
     result := average(add,count1)
     if result>=60 {
              a="成绩合格"
        }else {
              a="成绩不合格"
        }
     fmt.Println("平均成绩为",result,a)
      }