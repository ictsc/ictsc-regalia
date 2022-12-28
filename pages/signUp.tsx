import React, {useState} from "react";
import {useRouter} from "next/router";

import {SubmitHandler, useForm} from "react-hook-form";

import ICTSCNavBar from "../components/Navbar";
import {ICTSCSuccessAlert, ICTSCErrorAlert} from "../components/Alerts";
import {useApi} from "../hooks/api";

type Inputs = {
  name: string;
  password: string;
}

const SignUp = () => {
  const router = useRouter();
  const {user_group_id, invitation_code} = router.query;

  const {register, handleSubmit, formState: {errors}} = useForm<Inputs>()


  const {apiClient} = useApi();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null)
  // エラーメッセージ
  const [message, setMessage] = useState<string | null>(null)

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const response = await apiClient.post("users", {
      json: {
        ...data,
        user_group_id: user_group_id,
        invitation_code: invitation_code
      },
    })

    setStatus(response.status)

    if (!response.ok) {
      const message = await response.text()
      if (message.match(/Error 1062: Duplicate entry '\w+' for key 'name'/)) {
        setMessage("ユーザー名が重複しています。")
      }
      if (message.match(/Error:Field validation for 'UserGroupID' failed on the 'required' tag/)) {
        setMessage("無効なユーザーグループです。")
      }
      if (message.match(/Error:Field validation for 'UserGroupID' failed on the 'uuid' tag/)) {
        setMessage("無効なユーザーグループです。")
      }
      if (message.match(/Error:Field validation for 'InvitationCode' failed on the 'required' tag/)) {
        setMessage("無効な招待コードです。")
      }
    }

    if (response.status === 201) {
      await router.push("/login")
    }
  }


  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>ユーザー登録</h1>
        <form onSubmit={handleSubmit(onSubmit)}
              className={'form-control flex flex-col container-ictsc items-center'}>
          {status === 201 &&
            <ICTSCSuccessAlert className={'mb-8'} message={"ユーザー登録に成功しました！"}/>}
          {(status != null && status !== 201) &&
            <ICTSCErrorAlert className={'mb-8'} message={"エラーが発生しました"} subMessage={message ?? ''}/>}
          <input {...register('name', {required: true})}
                 type="text" placeholder="ユーザー名"
                 className="input input-bordered max-w-xs min-w-[312px]"/>
          <label className="label max-w-xs min-w-[312px]">
            {errors.name && <span className="label-text-alt text-error">ユーザー名を入力してください</span>}
          </label>
          <input {...register('password', {
            required: true,
            minLength: 8,
          })}
                 type="password" placeholder="パスワード"
                 className="input input-bordered max-w-xs min-w-[312px] mt-4"/>
          <label className="label max-w-xs min-w-[312px]">
            {errors.password?.type == "required" &&
              <span className="label-text-alt text-error">パスワードを入力して下さい</span>}
            {errors.password?.type == "minLength" &&
              <span className="label-text-alt text-error">パスワードは8文字以上である必要があります</span>}
          </label>
          <input type={"submit"} value={'登録'}
                 className="btn btn-primary mt-4 max-w-xs min-w-[312px]"/>
        </form>
      </>
  );
}

export default SignUp