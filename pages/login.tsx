import {useState} from "react";
import {useRouter} from "next/router";

import {SubmitHandler, useForm} from "react-hook-form";

import ICTSCNavBar from "../components/Navbar";
import {useAuth} from "../hooks/auth";
import {useApi} from "../hooks/api";


type Inputs = {
  name: string;
  password: string;
}

const Login = () => {
  const router = useRouter();

  const {register, handleSubmit, formState: {errors}} = useForm<Inputs>()

  const {apiClient} = useApi();
  const {mutate} = useAuth()

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null)

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const response = await apiClient.post("auth/signin", {
      json: data,
    })

    setStatus(response.status)

    if (response.status === 200) {
      await mutate()
      await router.push("/")
    }
  }

  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>ログイン</h1>
        <form onSubmit={handleSubmit(onSubmit)}
              className={'form-control flex flex-col container-ictsc items-center'}>
          {status === 200
              ? <LoginSuccessAlert/>
              : <LoginFailedAlert/>}
          <input {...register('name', {required: true})}
                 type="text" placeholder="ユーザー名"
                 className="input input-bordered max-w-xs min-w-[312px]"/>
          <label className="label max-w-xs min-w-[312px]">
            {errors.name && <span className="label-text-alt text-error">ユーザー名を入力してください</span>}
          </label>
          <input {...register('password', {required: true})}
                 type="password" placeholder="パスワード"
                 className="input input-bordered max-w-xs min-w-[312px] mt-4"/>
          <label className="label max-w-xs min-w-[312px]">
            {errors.password && <span className="label-text-alt text-error">パスワードを入力して下さい</span>}
          </label>
          <input type={"submit"} value={'ログイン'}
                 className="btn btn-primary mt-4 max-w-xs min-w-[312px]"/>
        </form>
      </>
  );
}

// ログインが成功した時に表示するアラート
const LoginSuccessAlert = () => {
  return (
      <div className="alert alert-success shadow-lg max-w-xs min-w-[312ppx] mb-8">
        <div>
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current flex-shrink-0 h-6 w-6" fill="none"
               viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span>ログインに成功しました！</span>
        </div>
      </div>
  )
}

// ログインが失敗した時に表示するアラート
const LoginFailedAlert = () => {
  return (
      <div className='alert alert-error shadow-lg max-w-xs min-w-[312ppx] mb-8'>
        <div>
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current flex-shrink-0 h-6 w-6" fill="none"
               viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span>ログインに失敗しました</span>
        </div>
      </div>
  )
}

export default Login