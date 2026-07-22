export default function PaginaIndex() {
    return (
        <div className="space-y-6">
            <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
                <h1 className="text-2xl font-bold text-gray-800">Página Principal</h1>
                <p className="text-gray-600 mt-2">
                    Esta vista está renderizada dentro del área central con scroll independiente.
                </p>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
                <h2 className="text-lg font-semibold text-gray-700 mb-4">Prueba de Scroll</h2>
                <div className="h-[1000px] bg-gradient-to-b from-blue-50 to-indigo-50 border-2 border-dashed border-blue-200 rounded-lg p-6 flex flex-col justify-between">
                    <p className="text-blue-600 font-medium">↑ Inicio del área scrolleable</p>
                    <p className="text-blue-600 font-medium text-center">Desplázate hacia abajo...</p>
                    <p className="text-blue-600 font-medium">↓ Fin del área scrolleable (El footer se queda abajo sin aplastarse)</p>
                </div>
            </div>
        </div>
    );
}