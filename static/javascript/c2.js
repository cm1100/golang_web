



document.getElementById('search_avail').addEventListener('click',()=>{
  //notify("welcome here","warning")
  let html = `
    <form id="check-avail-form" acttion="#" method="post" novalidate class="needs-validation">
      <div class="row">
        <div class="col">
          <div class="row" id="reservation-model">
            <div class="col">
              <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">

            </div>
            <div class="col">
              <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">

            </div>


          </div>
        </div>
      </div>
    </form>
  `

  let attention = promt()
  //attention.success({"msg":"am here",icon:'success'})


  attention.custom({msg:html,title:"Choose Your Dates",callback:(result)=>{


    let form = document.getElementById("check-avail-form")
    let formData = new FormData(form)
    formData.append("csrf_token","{{.CSRFToken}}")
    //console.log(formdata.keys()[2],formdata.values())

    fetch('/book-json',{
      method:"post",
      body:formData,
    }).then(response=>response.json).then(data=>{
      console.log(data)
    })

  }})
})



function promt(){

  let toast = function(c){

            const {
              msg= "",
              icon = "success",
              position = 'top-end',

            }=c;
            const Toast = Swal.mixin({
                toast: true,
                title:msg,
                position: position,
                icon : icon,
                showConfirmButton: false,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
                }
            })

            Toast.fire({})

  }

  let success = function(c){

    const {msg="",title="",footer=null,icon='success'}=c
    Swal.fire({
          icon: icon,
          title: title,
          text: msg,
          footer: footer
            })
  }



  async function custom(c){
    const {
      msg="",
      title="",
    }=c;
    const { value: result } = await Swal.fire({

  title: title,
  html:msg,
  backdrop:false,
  focusConfirm: false,
  showCancelButton:true,
  position : 'center',
  willOpen:()=>{
    const dates = document.getElementById('reservation-model');
    const rp = new DateRangePicker(dates,{
      format:'yyyy-mm-dd',
      showOnFocus: true,
      orientation:'top',

    })
  },
  preConfirm: () => {
    return [
      document.getElementById('start').value,
      document.getElementById('end').value
    ]
    },
  didOpen:()=>{
    document.getElementById('start').removeAttribute('disabled')
    document.getElementById('end').removeAttribute('disabled')
  },
  })

  if (result){
      if (result.dismiss !== Swal.DismissReason.cancle){
        if (result.value !== ""){
          if (c.callback !== undefined){
            c.callback(result)
          }
        }else{
          c.callback(false)
        }
      }else{
        c.callback(false)
      }
    }
  }

  return {
    toast:toast,
    success:success,
    custom:custom,
  }
}

