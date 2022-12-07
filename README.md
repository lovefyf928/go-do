# go-do

a mix many tools framework of microservice


main.go
```
server.Start("your config path", "your service global prefix", func(r *register.Register) {
		controller.DeviceController(r) // register your controller
	})
  ```
  
  
  DeviceController.go
```
func DeviceController(r *register.Register) {
	r.RegisterController("your controller prefix")

	r.RegisterHttpHandle("service path", enum.GET(request method), func(ctx context.Context, req interface{}) (interface{}, error) {
		tp, ok := req.(*params.TestParams)
		if ok {
			return deviceService.GetDeviceInfo(tp) // your logic code
		}
		return nil, errors.New(string(common.PARAMS_ERROR))
	}(handle), &params.TestParams{}(input params))
}
  ```
