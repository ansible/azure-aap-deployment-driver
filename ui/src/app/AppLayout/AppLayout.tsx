import * as React from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router';
import {
  Brand,
  Nav,
  NavList,
  NavItem,
  NavExpandable,
  Page,
  PageHeader,
  PageSidebar,
  SkipToContent,
  Modal, ModalVariant, Button
} from '@patternfly/react-core';
import logo from '../bgimages/Technology_icon-Red_Hat-Ansible_Automation_Platform-Standard-RGB.svg';
import { logout } from "../apis/auth";

interface IAppLayout {
  navigation: any[]
}
interface LoadingPropsType {
  spinnerAriaValueText: string;
  spinnerAriaLabelledBy?: string;
  spinnerAriaLabel?: string;
  isLoading: boolean;
}


const AppLayout = ({ navigation }: IAppLayout) => {
  const [isNavOpen, setIsNavOpen] = React.useState(true);
  const [isMobileView, setIsMobileView] = React.useState(true);
  const [isNavOpenMobile, setIsNavOpenMobile] = React.useState(false);
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [isPrimaryLoading, setIsPrimaryLoading] = React.useState<boolean>(false);

  const primaryLoadingProps = {} as LoadingPropsType;
  primaryLoadingProps.spinnerAriaValueText = 'Loading';
  primaryLoadingProps.spinnerAriaLabelledBy = 'primary-loading-button';
  primaryLoadingProps.isLoading = isPrimaryLoading;
  const navigate = useNavigate();
  const onNavToggleMobile = () => {
    setIsNavOpenMobile(!isNavOpenMobile);
  };
  const onNavToggle = () => {
    setIsNavOpen(!isNavOpen);
  };
  const onPageResize = (props: { mobileView: boolean; windowSize: number }) => {
    setIsMobileView(props.mobileView);
  };

  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
  };

  async function submitHandler(event) {
    event.preventDefault();
    setIsPrimaryLoading(!isPrimaryLoading);
    await logout();
    setTimeout(() => {
      setIsModalOpen(!isModalOpen);
      navigate("/login")
    }, 3000)
  }

  const Header = (
    <PageHeader
      logo={<Brand src={logo} alt="Red Hat Ansible Automation Platform Logo"/>}
      showNavToggle
      isNavOpen={isNavOpen}
      onNavToggle={isMobileView ? onNavToggleMobile : onNavToggle}
    />
  );

  const location = useLocation();

  const renderNavItem = (route: any, index: number) => {
    if (route.label !== "Logout") {
      return (
        <NavItem key={`${route.label}-${index}`} id={`${route.label}-${index}`} isActive={route.path === location.pathname} to={route.path}>
          {route.label}
        </NavItem>
      )
    }
    else {
      return (
        <NavItem key={`${route.label}-${index}`} id={`${route.label}-${index}`} isActive={route.path === location.pathname} onClick={handleModalToggle}>
          {route.label}
        </NavItem>
      )
    }
  };

  const renderNavGroup = (group: any, groupIndex: number) => (
    <NavExpandable
      key={`${group.label}-${groupIndex}`}
      id={`${group.label}-${groupIndex}`}
      title={group.label}
      isActive={group.routes.some((route) => route.path === location.pathname)}
    >
      {group.routes.map((route, idx) => route.label && renderNavItem(route, idx))}
    </NavExpandable>
  );

  const Navigation = (
    <Nav id="nav-primary-simple" theme="dark">
      <NavList id="nav-list-simple">
        {navigation.map(
          (route, idx) => route.label && (!route.routes ? renderNavItem(route, idx) : renderNavGroup(route, idx))
        )}
      </NavList>
    </Nav>
  );

  const Sidebar = (
    <PageSidebar
      theme="dark"
      nav={Navigation}
      isNavOpen={isMobileView ? isNavOpenMobile : isNavOpen} />
  );

  const pageId = 'primary-app-container';

  const PageSkipToContent = (
    <SkipToContent onClick={(event) => {
      event.preventDefault();
      const primaryContentContainer = document.getElementById(pageId);
      primaryContentContainer && primaryContentContainer.focus();
    }} href={`#${pageId}`}>
      Skip to Content
    </SkipToContent>
  );

  return (
    <Page
      mainContainerId={pageId}
      header={Header}
      sidebar={Sidebar}
      onPageResize={onPageResize}
      skipToContent={PageSkipToContent}
      className="pf-m-full-height" >
      <Outlet />
      <Modal
        variant={ModalVariant.small}
        title="Logout"
        isOpen={isModalOpen}
        onClose={handleModalToggle}
        actions={[
          <Button key="Logout" id='primary-loading-button' variant="primary" onClick={submitHandler} {...primaryLoadingProps}>
            {isPrimaryLoading ? 'Logging Out' : 'Confirm'}
          </Button>,
          <Button key="cancel" variant="link" onClick={handleModalToggle}>
            Cancel
          </Button>
        ]}
      >
        Are you sure you want Logout?
      </Modal>
    </Page>
  );
};

export { AppLayout };
