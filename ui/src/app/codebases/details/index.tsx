import CodeMap from '../../../components/CodeMap';
import { useParams } from 'react-router-dom';

const CodeMapPage = () => {
    const { codebaseId } = useParams();

    return (
        <CodeMap codebaseId={codebaseId} />
    );
}

export default CodeMapPage;