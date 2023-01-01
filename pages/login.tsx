import {useState} from "react";
import {useRouter} from "next/router";

import {SubmitHandler, useForm} from "react-hook-form";

import ICTSCNavBar from "../components/Navbar";
import {ICTSCErrorAlert, ICTSCSuccessAlert} from "../components/Alerts";
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
  // 送信中
  const [submitting, setSubmitting] = useState(false)

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    setSubmitting(true)
    const response = await apiClient.post("auth/signin", {
      json: data,
    })

    setSubmitting(false)
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
          {status === 200 &&
            <ICTSCSuccessAlert className={'mb-8'} message={'ログインに成功しました'}/>}
          {(status != null && status !== 200) &&
            <ICTSCErrorAlert className={'mb-8'} message={'ログインに失敗しました'}/>}
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
          <button type={"submit"}
                  className={`btn btn-primary mt-4 max-w-xs min-w-[312px] ${submitting && 'loading'}`}>
            ログイン
          </button>
        </form>
      </>
  );
}

export default Login