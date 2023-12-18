import * as React from 'react';
import Modal from '@mui/joy/Modal';

export default function ModalWrapper({
  open,
  onClose,
  children,
}: {
  open: boolean;
  onClose: () => void;
  children: React.ReactElement;
}) {
  return (
    <React.Fragment>
      <Modal open={open} onClose={onClose}>
        {children}
      </Modal>
    </React.Fragment>
  );
}
