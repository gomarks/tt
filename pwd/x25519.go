package main

import (
	"bytes"
	"context"
	x25519 "github.com/HirbodBehnam/EasyX25519"
	"golang.org/x/sync/semaphore"
	"log"
	"sync"
	"tt/common"
)

// x25519 是一种秘钥交换算法，基于椭圆曲线理论
// 采用定义在 GF(2^255 - 19)上的椭圆曲线：y^2 = x^3 + 48662 * x^2 + x，条件：P/Q 两点做直线，和曲线焦点为 R = P + Q，如果 P = Q，那么 nP = P1 + P2 …… + Pn
// 释义：
// A = aG
// B = bG
// A 发给 bob，计算 bA = b（aG） ？ 如果只是乘法运算，那么私钥 a 不就泄露了？
// B 发给 Alice，计算 aB = a（bG）
// 得到共享秘钥 K = abG

// 结论：x25519 很快：
const (
	ScalarSize = 32
	PointSize  = 32
)

func main() {

	common.ShowTime(100000, 0, func(sem *semaphore.Weighted, waitGroup *sync.WaitGroup, idx int) {
		if sem != nil {
			sem.Acquire(context.Background(), 1)
		}
		aliceKeyPair, err := x25519.NewX25519()
		if err != nil {
			log.Fatalln(err)
		}

		bobKeyPair, err := x25519.NewX25519()
		if err != nil {
			log.Fatalln(err)
		}

		aliceSwichKey, err := aliceKeyPair.GenerateSharedSecret(bobKeyPair.PublicKey)
		if err != nil {
			log.Fatalln(err)
		}

		bobSwichKey, err := bobKeyPair.GenerateSharedSecret(aliceKeyPair.PublicKey)
		if err != nil {
			log.Fatalln(err)
		}

		if !bytes.Equal(aliceSwichKey, bobSwichKey) {
			log.Printf("alice private key:%s public key : %s", aliceKeyPair.SecretKey, aliceKeyPair.PublicKey)
			log.Printf("bob private key:%s public key : %s", bobKeyPair.SecretKey, bobKeyPair.PublicKey)
			log.Printf("is same swith key : %v %v", aliceSwichKey, bobSwichKey)
		}
		log.Print(idx)
		if sem != nil {
			sem.Release(1)
		}
		waitGroup.Done()
	})

}
