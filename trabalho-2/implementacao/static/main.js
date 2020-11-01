const form = document.getElementById('firebase-form');
const blockname = document.getElementById('block-name');

const cleanInput = () => (blockname.value = '');

const saveBlock = async (name) => {
  try {
    const req = await fetch(
      'https://evening-cove-12029.herokuapp.com/add-block',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        mode: 'cors',
        body: JSON.stringify({ name }),
      }
    );
    const res = await req.text();
    alert(res);
  } catch (e) {
    console.error(e);
  }
};

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  const userInput = blockname.value;

  if (userInput.length === 0) {
    alert('Tamanho inválido.');
    return;
  }

  await saveBlock(userInput);
  cleanInput();
});
