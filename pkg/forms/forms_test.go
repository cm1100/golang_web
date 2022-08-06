package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)



func TestForm_Valid(t *testing.T){
	r , _ := http.NewRequest("POST","/wherevver",nil)



	form := New(r.PostForm)

	if !form.Valid(){
		t.Error("form should be valid here")
	}


}


func TestForm_Required(t *testing.T){

	r , _ := http.NewRequest("POST","/wherevver",nil)

	form := New(r.PostForm)

	form.Required("a","b","c")
	if form.Valid(){
		t.Error("form should not be valid here")
	}

	postData := url.Values{}

	postData.Add("a","a")
	postData.Add("b","b")
	postData.Add("c","c")


	r.PostForm =  postData

	form = New(r.PostForm)

	if !form.Valid(){
		t.Error("form should be valid after adding the postdata")
	}


}




func TestFormHas(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)

	form := New(r.PostForm)

	has:= form.Has("whatever",r)
	if has{
		t.Error("This should not exist")
	}

	postedData:= url.Values{}
	postedData.Add("a","abc")
	r.PostForm= postedData
	form = New(r.PostForm)
	has =form.Has("a",r)
	if !has{
		t.Error("this should exist")
	}
}


func TestForm_MinLength(t *testing.T){

	r := httptest.NewRequest("POST","/whatever",nil)

	form := New(r.PostForm)

	form.MinLength("x",10,r)
	if form.Valid(){
		t.Error("it should be valid")
	}


	postedData := url.Values{}

	postedData.Add("somefield","some value")

	form = New(postedData)

	form.MinLength("somefield",5,r)
	if !form.Valid(){
		t.Error("should work")
	}


	form.MinLength("somefield",100,r)
	if form.Valid(){
		t.Error("This should not work")
	}


}


func TestForm_Email(t *testing.T){

	r:= httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)


	form.IsEmail("x")
	if form.Valid(){
		t.Error("this should be valid")
	}
	postData := url.Values{}
	postData.Add("email","abc@gmail.com")
	form = New(postData)
	form.IsEmail("email")
	if !form.Valid(){
		t.Error("this should  be valid")
	}



	postData = url.Values{}
	postData.Add("email","abc")
	form = New(postData)
	form.IsEmail("email")
	if form.Valid(){
		t.Error("this should not be valid")
	}

}














