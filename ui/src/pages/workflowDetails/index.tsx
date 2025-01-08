import { useParams } from 'react-router-dom';

const CodeMapPage = () => {
    const { codebaseId } = useParams();

    return (
        <>Details of {codebaseId}</>
    );
}

export default CodeMapPage;