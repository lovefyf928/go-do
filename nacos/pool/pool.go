package pool

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net/rpc"
	"sync"
	"time"
)

type GrpcConn struct {
	l     *sync.Mutex
	Conn  *grpc.ClientConn
	timer *time.Timer
}

type RpcConn struct {
	l     *sync.Mutex
	Conn  *rpc.Client
	timer *time.Timer
}

type Pool struct {
	l *sync.Mutex
	//connPool              *sync.Map
	rpcConnPool           map[string][]*RpcConn
	grpcConnPool          map[string][]*GrpcConn
	maxSingleInstanceConn int
	maxSingleInstanceTime time.Duration
	grpcFn                func(serverName string) (conn *grpc.ClientConn, err error)
	rpcFn                 func(serverName string) (conn *rpc.Client, err error)
}

func NewGrpcPool(fn func(serverName string) (conn *grpc.ClientConn, err error), maxSingleInstanceConn int, maxSingleInstanceTime time.Duration) *Pool {
	if maxSingleInstanceConn <= 0 {
		panic("连接数必须大于0")
	}

	pool := &Pool{
		l:                     &sync.Mutex{},
		grpcFn:                fn,
		maxSingleInstanceConn: maxSingleInstanceConn,
		//connPool:              &sync.Map{},
		grpcConnPool:          map[string][]*GrpcConn{},
		maxSingleInstanceTime: maxSingleInstanceTime,
	}

	return pool
}

func NewRpcPool(fn func(serverName string) (conn *rpc.Client, err error), maxSingleInstanceConn int, maxSingleInstanceTime time.Duration) *Pool {
	if maxSingleInstanceConn <= 0 {
		panic("连接数必须大于0")
	}

	pool := &Pool{
		l:                     &sync.Mutex{},
		rpcFn:                 fn,
		maxSingleInstanceConn: maxSingleInstanceConn,
		//connPool:              &sync.Map{},
		rpcConnPool:           map[string][]*RpcConn{},
		maxSingleInstanceTime: maxSingleInstanceTime,
	}

	return pool
}

func (p *Pool) GetRpcConn(name string) (*RpcConn, error) {
	p.l.Lock()
	rpcConns := p.rpcConnPool[string(name)]
	if rpcConns != nil {
		//rpcConns, ok2 := arr.([]*RpcConn)
		//if ok2 {
		fmt.Println("len-------->", len(rpcConns))
		for i := range rpcConns {
			if rpcConns[i].l.TryLock() {
				p.l.Unlock()
				return rpcConns[i], nil
			}
		}
		if len(rpcConns) >= p.maxSingleInstanceConn {
			for true {
				for i := range rpcConns {
					if rpcConns[i].l.TryLock() {
						p.l.Unlock()
						return rpcConns[i], nil
					}
				}
			}
		} else {
			conn, err := p.rpcFn(name)
			if err != nil {
				p.l.Unlock()
				return nil, err
			}
			rpcConn := &RpcConn{
				l:     &sync.Mutex{},
				Conn:  conn,
				timer: time.NewTimer(p.maxSingleInstanceTime),
			}
			rpcConn.l.Lock()
			rpcConns = append(rpcConns, rpcConn)
			p.rpcConnPool[string(name)] = rpcConns
			//p.connPool.Store(name, rpcConns)
			p.l.Unlock()
			return rpcConn, nil
		}
		//}
	} else {
		conn, err := p.rpcFn(name)
		if err != nil {
			p.l.Unlock()
			return nil, err
		}
		rpcConn := &RpcConn{
			l:     &sync.Mutex{},
			Conn:  conn,
			timer: time.NewTimer(p.maxSingleInstanceTime),
		}
		rpcConn.l.Lock()
		rpcConns := []*RpcConn{rpcConn}
		p.rpcConnPool[string(name)] = rpcConns
		//p.connPool.Store(name, rpcConns)
		p.l.Unlock()
		return rpcConn, nil
	}
	p.l.Unlock()
	return nil, errors.New("内部错误")
}

func (p *Pool) GetGrpcConn(name string) (*GrpcConn, error) {
	p.l.Lock()
	grpcConns := p.grpcConnPool[string(name)]
	fmt.Println("len-------->", len(grpcConns))
	if grpcConns != nil {
		//grpcConns, ok2 := arr.([]*GrpcConn)
		//if ok2 {
		for i := range grpcConns {
			if grpcConns[i].l.TryLock() {
				p.l.Unlock()
				return grpcConns[i], nil
			}
		}
		if len(grpcConns) >= p.maxSingleInstanceConn {
			for true {
				for i := range grpcConns {
					if grpcConns[i].l.TryLock() {
						p.l.Unlock()
						return grpcConns[i], nil
					}
				}
			}
		} else {
			conn, err := p.grpcFn(name)
			if err != nil {
				p.l.Unlock()
				return nil, err
			}
			grpcConn := &GrpcConn{
				l:     &sync.Mutex{},
				Conn:  conn,
				timer: time.NewTimer(p.maxSingleInstanceTime),
			}
			grpcConn.l.Lock()
			grpcConns = append(grpcConns, grpcConn)
			p.grpcConnPool[string(name)] = grpcConns
			//p.connPool.Store(name, grpcConn)
			p.l.Unlock()
			return grpcConn, nil
		}
		//} else {
		//	//grpcc, ok1 := arr.([]*GrpcConn)
		//	//if ok1 {
		//	//	fmt.Println(grpcc)
		//	//} else {
		//	//	panic(arr)
		//	//}
		//	////panic(arr)
		//	p.l.Unlock()
		//	return nil, errors.New("123123123123123")
		//}
	} else {
		conn, err := p.grpcFn(name)
		if err != nil {
			p.l.Unlock()
			return nil, err
		}
		grpcConn := &GrpcConn{
			l:     &sync.Mutex{},
			Conn:  conn,
			timer: time.NewTimer(p.maxSingleInstanceTime),
		}
		grpcConn.l.Lock()
		grpcConns := []*GrpcConn{grpcConn}
		p.grpcConnPool[string(name)] = grpcConns
		//p.connPool.Store(name, grpcConns)
		p.l.Unlock()
		return grpcConn, nil
	}
	p.l.Unlock()
	return nil, errors.New("内部错误")
}

//func (p *Pool[T]) GetRpcConn(name enum.ServerName) (*RpcConn, error) {
//	p.l.Lock()
//	arr, ok := p.connPool.Load(name)
//	if ok {
//		rpcConns, ok2 := arr.([]*RpcConn)
//		if ok2 {
//			fmt.Println("len-------->", len(rpcConns))
//			for i := range rpcConns {
//				if rpcConns[i].l.TryLock() {
//					p.l.Unlock()
//					return rpcConns[i], nil
//				}
//			}
//			if len(rpcConns) >= p.maxSingleInstanceConn {
//				for true {
//					for i := range rpcConns {
//						if rpcConns[i].l.TryLock() {
//							p.l.Unlock()
//							return rpcConns[i], nil
//						}
//					}
//				}
//			} else {
//				conn, err := p.rpcFn(name)
//				if err != nil {
//					p.l.Unlock()
//					return nil, err
//				}
//				rpcConn := &RpcConn{
//					l:     &sync.Mutex{},
//					Conn:  conn,
//					timer: time.NewTimer(p.maxSingleInstanceTime),
//				}
//				rpcConn.l.Lock()
//				rpcConns = append(rpcConns, rpcConn)
//				p.connPool.Store(name, rpcConns)
//				p.l.Unlock()
//				return rpcConn, nil
//			}
//		}
//	} else {
//		conn, err := p.rpcFn(name)
//		if err != nil {
//			p.l.Unlock()
//			return nil, err
//		}
//		rpcConn := &RpcConn{
//			l:     &sync.Mutex{},
//			Conn:  conn,
//			timer: time.NewTimer(p.maxSingleInstanceTime),
//		}
//		rpcConn.l.Lock()
//		rpcConns := []*RpcConn{rpcConn}
//		p.connPool.Store(name, rpcConns)
//		p.l.Unlock()
//		return rpcConn, nil
//	}
//	p.l.Unlock()
//	return nil, errors.New("内部错误")
//}
//
//func (p *Pool[T]) GetGrpcConn(name enum.ServerName) (*GrpcConn, error) {
//	p.l.Lock()
//	arr, ok := p.connPool.Load(name)
//	if ok {
//		grpcConns, ok2 := arr.([]*GrpcConn)
//		if ok2 {
//			for i := range grpcConns {
//				if grpcConns[i].l.TryLock() {
//					p.l.Unlock()
//					return grpcConns[i], nil
//				}
//			}
//			if len(grpcConns) >= p.maxSingleInstanceConn {
//				for true {
//					for i := range grpcConns {
//						if grpcConns[i].l.TryLock() {
//							p.l.Unlock()
//							return grpcConns[i], nil
//						}
//					}
//				}
//			} else {
//				conn, err := p.grpcFn(name)
//				if err != nil {
//					p.l.Unlock()
//					return nil, err
//				}
//				grpcConn := &GrpcConn{
//					l:     &sync.Mutex{},
//					Conn:  conn,
//					timer: time.NewTimer(p.maxSingleInstanceTime),
//				}
//				grpcConn.l.Lock()
//				grpcConns = append(grpcConns, grpcConn)
//				p.connPool.Store(name, grpcConn)
//				p.l.Unlock()
//				return grpcConn, nil
//			}
//		} else {
//			//grpcc, ok1 := arr.([]*GrpcConn)
//			//if ok1 {
//			//	fmt.Println(grpcc)
//			//} else {
//			//	panic(arr)
//			//}
//			////panic(arr)
//			p.l.Unlock()
//			return nil, errors.New("123123123123123")
//		}
//	} else {
//		conn, err := p.grpcFn(name)
//		if err != nil {
//			p.l.Unlock()
//			return nil, err
//		}
//		grpcConn := &GrpcConn{
//			l:     &sync.Mutex{},
//			Conn:  conn,
//			timer: time.NewTimer(p.maxSingleInstanceTime),
//		}
//		grpcConn.l.Lock()
//		grpcConns := []*GrpcConn{grpcConn}
//		p.connPool.Store(name, grpcConns)
//		p.l.Unlock()
//		return grpcConn, nil
//	}
//	p.l.Unlock()
//	return nil, errors.New("内部错误")
//}

func (p *Pool) ReturnConn(instance interface{}) {
	rpcConn, ok := instance.(*RpcConn)
	if ok {
		rpcConn.timer.Reset(p.maxSingleInstanceTime)
		rpcConn.l.Unlock()
	}

	grpcConn, ok := instance.(*GrpcConn)
	if ok {
		grpcConn.timer.Reset(p.maxSingleInstanceTime)
		grpcConn.l.Unlock()
	}
}

func (p *Pool) DelConnGroup(groupName string) {

	delete(p.grpcConnPool, string(groupName))
	delete(p.rpcConnPool, string(groupName))

	//rpcConn, ok := instance.(*RpcConn)
	//if ok {
	//	rpcConn.timer.Reset(p.maxSingleInstanceTime)
	//	rpcConn.l.Unlock()
	//}
	//
	//grpcConn, ok := instance.(*GrpcConn)
	//if ok {
	//	grpcConn.timer.Reset(p.maxSingleInstanceTime)
	//	grpcConn.l.Unlock()
	//}
}

func (p *Pool) closeExpiredRpcConn(timer *time.Timer, rpcConn *RpcConn) {
	<-timer.C
	if !rpcConn.l.TryLock() {
		timer.Reset(p.maxSingleInstanceTime)
		return
	}
	rpcConn.Conn.Close()
}

func (p *Pool) closeExpiredGrpcConn(timer *time.Timer, grpcConn *GrpcConn) {
	<-timer.C
	if !grpcConn.l.TryLock() {
		timer.Reset(p.maxSingleInstanceTime)
		return
	}
	grpcConn.Conn.Close()
}
