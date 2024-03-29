import React from 'react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../apis/auth';
import { FixedUsernameLoginForm } from './FixedUsernameLoginForm';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import './LoginForm.css';

export function FormLogin() {
  const [loginMessage, setLoginMessage] = useState(
    'Please use the administrative credentials for Red Hat Ansible Automation Platform on Microsoft Azure.'
  );
  const navigate = useNavigate();

  async function loginHandler(uid: string, pwd: string) {
    try {
      const response = await login({ uid: uid, pwd: pwd });
      if ('unavailable' === response.status) {
	setLoginMessage('Sorry, we are having trouble signing you in as your deployment may have failed. Please return to Microsoft Azure Portal and attempt to redeploy Ansible Automation Platform.')
        setShowHelperText(true)
        setIsLoginButtonDisabled(true)
        setPasswordDisabled(true)
      }
      else if ('error' in response && response.error) {
        setLoginMessage('Incorrect Password');
        setIsValidPassword(false)
        setShowHelperText(true);
      } else {
        setIsValidPassword(true)
        navigate('/rhlogin');
      }
    } catch (err: any) {
      console.log('Got exception from logging in.', err);
      setIsValidPassword(false)
      if ('message' in err && err.message) {
        setLoginMessage(err.message);
      } else {
        setLoginMessage('Login failed. Try again...');
      }
    }
  }

  const [showHelperText, setShowHelperText] = React.useState(false);
  const [password, setPassword] = React.useState('');
  const [isValidPassword, setIsValidPassword] = React.useState(true);
  const [isLoginButtonDisabled, setIsLoginButtonDisabled] = React.useState(false);
  const [passwordDisabled, setPasswordDisabled] = React.useState(false);

  const handlePasswordChange = (value: string) => {
    setPassword(value);
  };

  const onLoginButtonClick = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    setShowHelperText(!password);
    loginHandler('admin', password);
  };

  return (
    <FixedUsernameLoginForm
      showHelperText={showHelperText}
      helperText={loginMessage}
      helperTextIcon={<ExclamationCircleIcon />}
      usernameValue={'admin'}
      passwordLabel="Password"
      passwordValue={password}
      isShowPasswordEnabled
      onChangePassword={handlePasswordChange}
      isValidPassword={isValidPassword}
      onLoginButtonClick={onLoginButtonClick}
      loginButtonLabel="Log in"
      usernameDisabled={true}
      isLoginButtonDisabled={isLoginButtonDisabled}
      passwordDisabled={passwordDisabled}
    />
  );
}
