function notify(msg,msgtype){
  notie.alert({
  type: msgtype, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
  text: msg,
  //stay: Boolean, // optional, default = false
  //time: Number, // optional, default = 3, minimum = 1,
  //position:  // optional, default = 'top', enum: ['top', 'bottom']
})
}


function notifyModal(title,text,icon,confirmButtonText,footer=null){

  Swal.fire({
    icon: icon ,
    title: title,
    text: text,
    confirmButtonText:confirmButtonText,
    footer: footer,
  })

}

document.getElementById('btn1').addEventListener('click',()=>{
  //notify("welcome here","warning")
  let html = `
    <form id="check-avail" acttion="#" method="post" novalidate class="needs-validation">
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
  attention.custom({msg:html,title:"Choose Your Dates"})
})

//prompt is javascript module for 
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
    const { value: formValues } = await Swal.fire({
  title: title,
  html:msg,
  backdrop:false,
  focusConfirm: false,
  showCancelButton:true,
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

  if (formValues) {
    Swal.fire(JSON.stringify(formValues))
    }
    }

  return {
    toast:toast,
    success:success,
    custom:custom,
  }
}


















