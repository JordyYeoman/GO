package main

import (
	"fmt"
	"math/rand"
	"time"
)

var imagePaths = [12]string{
	"https://images.unsplash.com/photo-1607627000458-210e8d2bdb1d?q=80&w=2949&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1566932769119-7a1fb6d7ce23?q=80&w=2970&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1579156618335-f6245e05236a?q=80&w=2960&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1616500156885-e51d834cab8e?q=80&w=2970&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1579156412503-f22426cc6386?q=80&w=2973&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1593013820725-ca0b6076576f?q=80&w=2970&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://plus.unsplash.com/premium_photo-1664304851973-5cb49d6a2230?q=80&w=2993&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1595210382266-2d0077c1f541?q=80&w=2970&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1517339763538-9bb6c388613a?q=80&w=2048&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1599409091912-88526846d833?q=80&w=2970&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://images.unsplash.com/photo-1595883696983-20e83efbd702?q=80&w=2974&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	"https://wallpapers-clan.com/wp-content/uploads/2023/11/aesthetic-iron-man-marvel-desktop-wallpaper-preview.jpg",
}

func main() {
	// Set the seed value for the random number generator
	rand.Seed(time.Now().UnixNano())

	randImg := imagePaths[rand.Intn(12)]
	fmt.Println(randImg)
	//for i := 0; i < 10000; i++ {
	//	randImg := imagePaths[rand.Intn(12)]
	//	//fmt.Println("do it")
	//	if randImg == "https://wallpapers-clan.com/wp-content/uploads/2023/11/aesthetic-iron-man-marvel-desktop-wallpaper-preview.jpg" {
	//		fmt.Println("We made it!!")
	//		fmt.Println(i)
	//		fmt.Println(randImg)
	//		break
	//	}
	//}
}
