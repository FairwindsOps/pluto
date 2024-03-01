Name:           pluto
Version:        5.19.0
Release:        1%{?dist}
Summary:        Fairwinds Kubernetes Configuration Validator

License:        Apache-2.0
URL:            https://github.com/FairwindsOps/pluto

Source0:        https://github.com/FairwindsOps/pluto/archive/v%{version}.tar.gz

BuildRequires:  make, git, go, wget, go
BuildRequires:  golang >= 1.16

%description
Pluto is a tool for validating Kubernetes configuration files.

%prep
%autosetup -n %{name}-%{version}

%build
export PATH=$PWD/go/bin:$PATH
go version
make %{?_smp_mflags} build buildid

%install
mkdir -p %{buildroot}/usr/bin
install -D -m 0755 %{_builddir}/pluto-%{version}/pluto %{buildroot}/usr/bin/pluto

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Fri Mar 01 2024 Emanuele Ciurleo <emanuele@ciurleo.com> - 5.19.0-1
- Initial package release