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
      if ('error' in response && response.error) {
        setLoginMessage('Incorrect Password');
        setShowHelperText(true);
      } else {
        navigate('/');
      }
    } catch (err: any) {
      console.log('Got exception from logging in.', err);
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

  const handlePasswordChange = (value: string) => {
    setPassword(value);
  };

  const onLoginButtonClick = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    setIsValidPassword(!!password);
    setShowHelperText(!'admin' || !password);
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
    />
  );
}
